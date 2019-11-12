package main

import (
	"net/http"
	"encoding/json"
	"github.com/gin-gonic/gin"
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

func checkVal() gin.HandlerFunc {
	return func (c *gin.Context) {
		idx, err := strconv.Atoi(c.Param("num"))
		if err != nil || idx < 0 {
			c.AbortWithStatusJSON(400, gin.H{
				"error": "Expected number greater or equal to 0",
			})
		}
		c.Next()
	}
}

func getAllVals(c *gin.Context) {
	var resp Resp
	res, err := http.Get("http://db_worker:3200/all")
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
