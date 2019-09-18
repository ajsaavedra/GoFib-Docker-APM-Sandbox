package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"strconv"
	"os"
)

var rdb *redis.Client

func main() {
	router := gin.New()

	setRedisClient()

	api := router.Group("/api")
	{
		api.GET("/fib/:memo/:num", checkVal(), publishIndex)
		api.GET("/all", getAllVals)
		api.POST("/fib/:num", postFibVal)
		api.DELETE("/fib/:num", deleteFibVal)
		api.DELETE("/all", deleteAllVals)
	}

	router.Run(":3100")
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func setRedisClient() {
	addr := os.Getenv("REDIS_HOST") + os.Getenv("REDIS_PORT")
	rdb = redis.NewClient(&redis.Options{
		Addr: addr,
		Password: "",           
	})

	_, err := rdb.Ping().Result()
	handleErr(err)
}

func publishIndex(c *gin.Context) {
	var channel string
	isMemo, _ := strconv.ParseBool(c.Param("memo"))
	if isMemo {
		channel = "memo-channel"
	} else {
		channel = "iter-channel"
	}

	err := rdb.Publish(channel, c.Param("num")).Err()
	handleErr(err)
}