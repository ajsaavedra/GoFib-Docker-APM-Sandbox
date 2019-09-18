package main

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func checkVal() gin.HandlerFunc {
	return func (c *gin.Context) {
		idx, err := strconv.Atoi(c.Param("num"))
		if err != nil || idx < 0 {
			c.JSON(400, gin.H{
				"error": "Expected number greater or equal to 0",
			})
		}
		c.Next()
	}
}

func getAllVals(c *gin.Context) {

}

func postFibVal(c *gin.Context) {

}

func deleteFibVal(c *gin.Context) {

}

func deleteAllVals(c *gin.Context) {

}
