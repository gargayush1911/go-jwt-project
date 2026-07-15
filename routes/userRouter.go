package routes

import(
	controller "go-jwt-project/controllers"
	"go-jwt-project/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingroutes *gin.Engine){
	incomingroutes.Use(middleware.Authentication())
	incomingroutes.GET("/users",controller.Getusers())
	incomingroutes.GET("/users/:user_id",controller.Getuser())
}