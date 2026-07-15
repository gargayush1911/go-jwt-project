package routes

import (
	controller "go-jwt-project/controllers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes( incomingroutes *gin.Engine){
	incomingroutes.POST("users/signup", controller.Signup())
	incomingroutes.POST("users/login", controller.Login())

}