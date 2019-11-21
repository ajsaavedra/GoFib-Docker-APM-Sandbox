package main

import (
	"net/http"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
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

type Msg struct {
	Body string `json:"payload"`
}

type Value struct {
	Num int `json:"value"`
}

func newValue() Value {
	return Value{ -1 }
}

func checkVal(c *gin.Context) {
	ctx := c.Request.Context()
	span, _ := tracer.StartSpanFromContext(ctx, "checkVal")

	value := newValue()
	json.NewDecoder(c.Request.Body).Decode(&value)

	if value.Num < 0 {
		c.AbortWithStatusJSON(400, gin.H{
			"error": "Bad request. Expected value field greater or equal to 0",
		})
	}

	c.Set("value", value.Num)
	span.Finish()
	c.Next()
}

func getVal(c *gin.Context) {
	url := fmt.Sprintf("http://db_worker:3200/fib/%s", c.Param("num"))
	httpClient := &http.Client{}
	httpReq, _ := http.NewRequest("GET", url, nil)

	var resp Fib
	res, err := httpClient.Do(httpReq)
	abortDBCall(err, c)

	defer res.Body.Close()

	if res.StatusCode == 404 {
		c.AbortWithStatusJSON(404, gin.H{
			"value": "Not found",
		})
		return
	}

	json.NewDecoder(res.Body).Decode(&resp)
	c.JSON(200, resp)
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
	var resp Msg
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