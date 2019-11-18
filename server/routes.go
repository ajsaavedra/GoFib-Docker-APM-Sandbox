package main

import (
	"net/http"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"strconv"
	"fmt"
)

type Fib struct {
	Idx int			`json:"idx"`
	Fib string		`json:"fib"`
	Elapsed string	`json:"elapsed"`
}

type Resp struct {
	Body []Fib `json:"payload"`
}

type Status struct {
	Body string `json:"payload"`
}

func checkVal(c *gin.Context) {
	ctx := c.Request.Context()
	span, _ := tracer.StartSpanFromContext(ctx, "checkVal")
	idx, err := strconv.Atoi(c.Param("num"))
	if err != nil || idx < 0 {
		c.AbortWithStatusJSON(400, gin.H{
			"error": "Expected number greater or equal to 0",
		})
	}
	span.Finish()
	c.Next()
}

func getAllVals(c *gin.Context) {
	httpClient := &http.Client{}
	httpReq, _ := http.NewRequest("GET", "http://db_worker:3200/all", nil)

	var resp Resp
	res, err := httpClient.Do(httpReq)
	abortDBCall(err, c)

	defer res.Body.Close()
	json.NewDecoder(res.Body).Decode(&resp)

	c.JSON(200, resp.Body);
}

func deleteFibVal(c *gin.Context) {
	client := &http.Client{}
	url := fmt.Sprintf("http://db_worker:3200/%s", c.Param("num"))
	var resp Status
	req, err := http.NewRequest("DELETE", url, nil)
	abortDBCall(err, c)

	res, err := client.Do(req)

	abortDBCall(err, c)

	defer res.Body.Close()
	json.NewDecoder(res.Body).Decode(&resp)

	c.JSON(200, resp.Body);
}

func abortDBCall(err error, c *gin.Context) {
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"error": "Could not query DB",
		})
	}
}