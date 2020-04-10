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

	// Router get info user
	r.GET("/get_info_user", models.GetInfoUser)

	// Router update info user (password, email)
	r.PATCH("/update_info_user", models.UpdateInfoUser)

	// Router insert user
	r.POST("/insert_user", models.InsertUser)

	// Router delete user
	r.DELETE("/delete_user", models.DeleteUser)

	r.Run()
}
