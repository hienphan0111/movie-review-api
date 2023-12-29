package routes

import (
	"github.com/hienphan0111/movie-review-api/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine) {
	router.POST("user/signup", controllers.Signup())
}
