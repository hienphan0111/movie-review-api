package main

import (
	"os"

	"github.com/hienphan0111/movie-review-api/database"
	"github.com/hienphan0111/movie-review-api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	router := gin.Default()

	//run database
	database.StartDB()

	//Log events
	router.Use(gin.Logger())

	routes.AuthRoutes(router)

	router.GET("/api", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"success": "Welcome to Movie Review API",
		})
	})

	router.Run(":" + port)
}
