package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-redis/redis"
	"strconv"
	"os"
)

var db *sql.DB
var dberr error
var rdb *redis.Client

func main() {
	router := gin.New()

	db, dberr = sql.Open("mysql", "root:password@/gofib")
	defer db.Close()

	setRedisClient()
	router.Run(":3200")
}

func setRedisClient() {
	addr := os.Getenv("REDIS_HOST") + os.Getenv("REDIS_PORT")
	rdb = redis.NewClient(&redis.Options{
		Addr: addr,
		Password: "",             
	})

	_, err := rdb.Ping().Result()
	if err != nil {
		panic(err)
	}

	subscribe()
}

func subscribe() {
	pubsub := rdb.Subscribe("memo-channel", "iter-channel")
	_, err := pubsub.Receive()
	if err != nil {
		panic(err)
	}

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

func insertFibValue(idx int) {
	fib := memoFib(idx, map[int]int{ 0:0, 1:1 })
	stmt, _ := db.Prepare("INSERT INTO sequences(idx, fib) VALUES (?, ?)")
	stmt.Exec(idx, fib)
}
