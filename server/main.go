package main

import (
	gintrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gin-gonic/gin"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/go-redis/redis"
	"strconv"
	"fmt"
	"os"
)

var rdb *redis.Client

func main() {
	debugMode, _ := strconv.ParseBool(os.Getenv("DEBUG_MODE"))

	tracer.Start(
		tracer.WithAgentAddr("datadog-agent:8126"),
		tracer.WithServiceName("api-service"),
		tracer.WithDebugMode(debugMode),
	)
	defer tracer.Stop()

	router := gin.New()
	router.Use(cors.Default())
	router.Use(gintrace.Middleware("go-fib"))

	setRedisClient()

	api := router.Group("/api")
	{
		api.GET("/fib/:num", checkVal, publishIndex)
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