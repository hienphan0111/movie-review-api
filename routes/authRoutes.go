package routes

import (
	"movie-review-api/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine) {
	router.POST("user/signup", controllers.Signup())
}
