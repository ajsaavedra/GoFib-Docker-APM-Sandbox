package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	sqltrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	log "github.com/sirupsen/logrus"
	"github.com/go-redis/redis"
	"net/http"
	"strconv"
	"context"
	"time"
	"fmt"
	"os"
)

var db *sql.DB
var dberr error
var rdb *redis.Client

type Fib struct {
	Idx int			`json:"idx"`
	Fib string		`json:"fib"`
	Elapsed string	`json:"elapsed"`
}

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	router := gin.Default()

	debugMode, _ := strconv.ParseBool(os.Getenv("DEBUG_MODE"))

	tracer.Start(
		tracer.WithAgentAddr("datadog-agent:8126"),
		tracer.WithServiceName("db-service"),
		tracer.WithDebugMode(debugMode),
	)
	defer tracer.Stop()

	sqltrace.Register("mysql", mysql.MySQLDriver{})
	cnxn := fmt.Sprintf("%s:%s@tcp(mysql)/%s",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PW"),
		os.Getenv("MYSQL_DB"),
	)
	db, dberr = sqltrace.Open("mysql", cnxn)
	defer db.Close()

	rdb = setRedisClient()
	go subscribe()

	router.GET("/all", getAllValues)
	router.GET("/fib/:num", getValue)
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
		go insertFibValue(msg.Payload)
	}
}

func getSpanFromContext(resourceName string) (tracer.Span, context.Context) {
	span, ctx := tracer.StartSpanFromContext(context.Background(), "fib-query",
		tracer.SpanType(ext.SpanTypeSQL),
		tracer.ServiceName("db-service"),
		tracer.ResourceName(resourceName),
	)

	return span, ctx
}

func getFibValue(msg string) (int, string, string) {
	span, ctx := getSpanFromContext("getFibValue")

	var idx int
	var fib, elapsed string
	num, _ := strconv.Atoi(msg)

	err := db.QueryRowContext(ctx, "SELECT * FROM sequences WHERE idx = ?", num).Scan(&idx, &fib, &elapsed)

	if err != nil {
		if err == sql.ErrNoRows {
			traceID := span.Context().TraceID()
			spanID := span.Context().SpanID()
			log.WithFields(log.Fields{"index": idx, "dd.trace_id": traceID, "dd.span_id": spanID}).Info("Value not found in DB")
		}
	}

	span.Finish(tracer.WithError(err))
	return idx, fib, elapsed
}

func getValue(c *gin.Context) {
	span, _ := getSpanFromContext("getValue")
	idx, fib, elapsed := getFibValue(c.Param("num"))
	if elapsed == "" {
		c.Writer.WriteHeader(http.StatusNotFound)
	} else {	
		value := Fib{idx, fib, elapsed}
		c.JSON(200, value)
	}
	span.Finish()
}

func getAllValues(c *gin.Context) {
	span, ctx := getSpanFromContext("getAllValues")
	rows, err := db.QueryContext(ctx, "SELECT idx, fib, elapsed FROM sequences")
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

	span.Finish(tracer.WithError(err))

	c.JSON(200, gin.H{
		"payload": values,
	});
}

func insertFibValue(msg string) {
	span, ctx := getSpanFromContext("insertFibValue")

	_, _, elapsedFound := getFibValue(msg)

	if elapsedFound != "" {
		return
	}

	idx, _ := strconv.Atoi(msg)
	start := time.Now()
	fib := memoFib(idx, map[int]int{ 0:0, 1:1 }, span)
	elapsed := time.Since(start).String()
	stmt, err := db.PrepareContext(ctx, "INSERT INTO sequences(idx, fib, elapsed) VALUES (?, ?, ?)")
	stmt.ExecContext(ctx, idx, fib, elapsed)

	traceID := span.Context().TraceID()
    spanID := span.Context().SpanID()
	log.WithFields(log.Fields{"index": idx, "value": fib, "dd.trace_id": traceID, "dd.span_id": spanID}).Info("Inserting calculated fib value")

	span.Finish(tracer.WithError(err))
}

func deleteFibValue(c *gin.Context) {
	span, ctx := getSpanFromContext("deleteFibValue")
	idx := c.Param("num")

	stmt, err := db.PrepareContext(ctx, "DELETE FROM sequences WHERE idx = ?")
	stmt.ExecContext(ctx, idx)

	c.JSON(200, gin.H{
		"payload": idx,
	})

	traceID := span.Context().TraceID()
    spanID := span.Context().SpanID()
	log.WithFields(log.Fields{"index": idx, "dd.trace_id": traceID, "dd.span_id": spanID}).Info("Deleted fib value from DB")

	span.Finish(tracer.WithError(err))
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
