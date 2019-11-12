package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-redis/redis"
	"encoding/json"
	"strconv"
	"time"
	"fmt"
	"os"
)

var db *sql.DB
var dberr error
var rdb *redis.Client
var rdbPub *redis.Client

func main() {
	router := gin.Default()
	cnxn := fmt.Sprintf("%s:%s@tcp(mysql)/%s",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PW"),
		os.Getenv("MYSQL_DB"),
	)
	db, dberr = sql.Open("mysql", cnxn)
	defer db.Close()

	rdb, rdbPub = setRedisClient(), setRedisClient()
	go subscribe()
	router.GET("/all", getAllValues)
	router.DELETE("/:num", deleteFibValue)
	router.Run(":3200")
}

func setRedisClient() *redis.Client {
	addr := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	client := redis.NewClient(&redis.Options{
		Addr: addr,
		Password: "",
	})

	_, err := client.Ping().Result()
	handleErr(err)
	return client
}

func subscribe() {
	pubsub := rdb.Subscribe("memo-channel")
	_, err := pubsub.Receive()
	handleErr(err)

	for msg := range pubsub.Channel() {
		go getFibValue(msg.Payload)
	}
}

func getFibValue(msg string) {
	var idx int
	var fib, elapsed string
	num, _ := strconv.Atoi(msg)

	err := db.QueryRow("SELECT * FROM sequences WHERE idx = ?", num).Scan(&idx, &fib, &elapsed)

	if err != nil {
		if err == sql.ErrNoRows {
			insertFibValue(num)
		}
	}
}

func getAllValues(c *gin.Context) {
	type Fib struct {
		Idx int			`json:"idx"`
		Fib string		`json:"fib"`
		Elapsed string	`json:"elapsed"`
	}

	rows, err := db.Query("SELECT idx, fib, elapsed FROM sequences")
	defer rows.Close()
	
	if err != nil {
		c.AbortWithStatus(500)
	}

	var values []Fib

	for rows.Next() {
		var idx int
		var fib, elapsed string
		rows.Scan(&idx, &fib, &elapsed)

		values = append(values, Fib{idx, fib, elapsed})
	}

	c.JSON(200, gin.H{
		"payload": values,
	});
}

func insertFibValue(idx int) {
	start := time.Now()
	fib := memoFib(idx, map[int]int{ 0:0, 1:1 })
	elapsed := time.Since(start).String()
	stmt, _ := db.Prepare("INSERT INTO sequences(idx, fib, elapsed) VALUES (?, ?, ?)")
	stmt.Exec(idx, fib, elapsed)
	emitFib(idx, fib, elapsed)
}

func deleteFibValue(c *gin.Context) {
	stmt, _ := db.Prepare("DELETE FROM sequences WHERE idx = ?")
	stmt.Exec(c.Param("num"))
	c.JSON(200, gin.H{
		"payload": c.Param("num"),
	})
}

func emitFib(idx int, fib int, elapsed string) {
	type data struct {
		Idx int
		Fib int
		Elapsed string
	}
	payload := &data{idx, fib, elapsed}
	message, _ := json.Marshal(payload)
	err := rdbPub.Publish("emit-channel", message).Err()
	handleErr(err)
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
