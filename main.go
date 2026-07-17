package main

import(
	routes "go-jwt-project/routes"
	"os"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
)

func main(){
	err := godotenv.Load(".env")
	if err!=nil{
		log.Fatal("error loading .env file")
	}
	port := os.Getenv("PORT")
	if port == ""{
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())

	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	router.GET("/api-1",func(c *gin.Context){
		c.JSON(200,gin.H{"success":"Access Granted For api-1"})
	})

	router.GET("/api-2",func(c *gin.Context){
		c.JSON(200,gin.H{"success":"Access Granted To api-2"})
	})

	router.Run(":"+port)


}