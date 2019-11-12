package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/go-redis/redis"
	"fmt"
	"os"
)

var rdb *redis.Client

func main() {
	router := gin.New()
	router.Use(cors.Default())

	setRedisClient()

	api := router.Group("/api")
	{
		api.GET("/fib/:num", checkVal(), publishIndex)
		api.GET("/all", getAllVals)
		api.DELETE("/:num", deleteFibVal)
	}

	router.Run(":3100")
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func setRedisClient() {
	addr := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	rdb = redis.NewClient(&redis.Options{
		Addr: addr,
		Password: "",           
	})

	_, err := rdb.Ping().Result()
	handleErr(err)
}

func publishIndex(c *gin.Context) {
	err := rdb.Publish("memo-channel", c.Param("num")).Err()
	handleErr(err)
	c.JSON(200, "done")
}