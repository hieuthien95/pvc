package main

import (
	"demo/database/models"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {

	initRouter()
}

func initRouter() {

	r := gin.Default()

	// Router test
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Router get info user
	r.GET("/get_info_user", models.GetInfoUser)

	r.Run()
}
