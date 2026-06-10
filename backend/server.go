package main

import (
	"backend/backend/controller"
	"backend/backend/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
			"http://localhost:3000",
		},
		AllowMethods: []string{"GET", "POST", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept"},
	}))

	router.Static("/screenshots", "./screenshots")

	router.POST("/report", controller.GetWebsiteReport)
	router.POST("/send-report", service.SendReport)

	router.Run(":8086")
}
