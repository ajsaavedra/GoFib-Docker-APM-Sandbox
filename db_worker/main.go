package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-redis/redis"
	"strconv"
	"fmt"
	"os"
)

var db *sql.DB
var dberr error
var rdb *redis.Client
var rdbPub *redis.Client

func main() {
	router := gin.Default()

	cnxn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PW"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_DB"),
	)
	db, dberr = sql.Open("mysql", cnxn)
	defer db.Close()

	rdb, rdbPub = setRedisClient(), setRedisClient()
	router.GET("/all", getAllValues)
	router.Run(":3200")

	go subscribe()
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
	pubsub := rdb.Subscribe("memo-channel", "iter-channel")
	_, err := pubsub.Receive()
	handleErr(err)

	for msg := range pubsub.Channel() {
		switch msg.Channel {
		case "memo-channel":
			go getFibValue(msg.Payload)
		case "iter-channel":
		}
	}
}

func getFibValue(msg string) {
	var idx, fib int
	num, _ := strconv.Atoi(msg)

	err := db.QueryRow("SELECT * FROM sequences WHERE idx = ?", num).Scan(&idx, &fib)

	if err != nil {
		insertFibValue(num)
	}
}

func getAllValues(c *gin.Context) {
	type Fib struct {
		Idx int	 `json:"idx"`
		Fib string `json:"fib"`
	}

	rows, err := db.Query("SELECT idx, fib FROM sequences")
	defer rows.Close()
	
	if err != nil {
		c.AbortWithStatus(500)
	}

	var values []Fib

	for rows.Next() {
		var idx int
		var fib string
		rows.Scan(&idx, &fib)

		values = append(values, Fib{idx, fib})
	}

	c.JSON(200, gin.H{
		"payload": values,
	});
}

func insertFibValue(idx int) {
	fib := memoFib(idx, map[int]int{ 0:0, 1:1 })
	stmt, _ := db.Prepare("INSERT INTO sequences(idx, fib) VALUES (?, ?)")
	stmt.Exec(idx, fib)
	emitFib(fib)
}

func emitFib(fib int) {
	err := rdbPub.Publish("emit-channel", fib).Err()
	handleErr(err)
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
