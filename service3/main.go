package main

import (
	"fmt"
	"inventory/configs"
	"inventory/controllers"
	"inventory/models"
	"inventory/repositories"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {
	configs.LoadEnv()
}

func main() {
	time.Sleep(1 * time.Minute)
	var connection = configs.Connect()
	connection.AutoMigrate(&models.Post{})

	var repository = repositories.PostRepository{Storage: connection}
	var drone = controllers.Post{Repository: &repository}
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	//r := gin.Default()

	r.POST("/posts", drone.Create)
	r.PUT("/posts/:id", drone.Update)
	r.DELETE("/posts/:id", drone.Delete)
	r.GET("/posts/:id", drone.Get)
	r.GET("/posts", drone.GetAll)

	r.Run(fmt.Sprintf(":%d", configs.Env.AppPort))
}
