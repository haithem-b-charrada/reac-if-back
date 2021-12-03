package main

import (
	"fmt"
	"inventory/configs"
	"inventory/controllers"
	"inventory/models"
	"inventory/repositories"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	configs.LoadEnv()
}

func main() {
	time.Sleep(1 * time.Minute)
	var connection = configs.Connect()
	connection.AutoMigrate(&models.User{})

	var repository = repositories.UserRepository{Storage: connection}
	var user = controllers.User{Repository: &repository}
	var login = controllers.Login{Repository: &repository}

	r := gin.Default()

	r.POST("/users", user.Create)
	r.PUT("/users/:id", user.Update)
	r.DELETE("/users/:id", user.Delete)
	r.GET("/users/:id", user.Get)
	r.GET("/users", user.GetAll)
	r.PUT("/password", user.ChangePassword)
	r.POST("/reset-password", user.ResetPasswordRequest)
	r.POST("/login", login.Auth)
	r.POST("/validate", login.Validate)
	r.POST("/set-new-password", user.SetNewPassword)
	r.POST("/seed", func(c *gin.Context) {
		hash, _ := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.MinCost)
		repository.Save(&models.User{
			Email:    "haithem.bencharrada@sastec-group.com",
			Password: string(hash),
		})

		c.JSON(http.StatusOK, gin.H{
			"message": "Seeding succuess",
		})
	})

	r.Run(fmt.Sprintf(":%d", configs.Env.AppPort))
}
