package main

import (
	"net/http"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"strconv"
)

type Fib struct {
	Idx int	   `json:"idx"`
	Fib string `json:"fib"`
}

type Resp struct {
	Body []Fib `json:"payload"`
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
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"error": "Could not query table",
		})
	}

	defer res.Body.Close()
	json.NewDecoder(res.Body).Decode(&resp)

	c.JSON(200, resp.Body);
}

func postFibVal(c *gin.Context) {

}

func deleteFibVal(c *gin.Context) {

}

func deleteAllVals(c *gin.Context) {

}
