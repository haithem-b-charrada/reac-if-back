package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"inventory/models"

	"inventory/repositories"

	"github.com/gin-gonic/gin"
)

type Post struct {
	Repository *repositories.PostRepository
}

func (d Post) Create(c *gin.Context) {
	var model = models.Post{}

	c.BindJSON(&model)

	d.Repository.Save(&model)

	c.JSON(http.StatusOK, model)
}

func (d Post) Update(c *gin.Context) {
	var model = models.Post{}
	//var request = models.Post{}

	model.ID, _ = strconv.Atoi(c.Param("id"))

	d.Repository.Find(&model)
	fmt.Println(model)
	if model.ID < 1 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "user not found",
			"status":  http.StatusNotFound,
		})
		return
	}

	c.BindJSON(&model)
	model.UpdatedAt = time.Now()

	d.Repository.Save(&model)

	c.JSON(http.StatusOK, model)
}

func (d Post) Delete(c *gin.Context) {
	var model = models.Post{}

	model.ID, _ = strconv.Atoi(c.Param("id"))

	d.Repository.Find(&model)
	if model.ID < 1 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "drone not found",
			"status":  http.StatusNotFound,
		})

		return
	}

	d.Repository.Remove(&model)

	c.JSON(http.StatusNoContent, "")
}

func (d Post) GetAll(c *gin.Context) {
	var models = []models.Post{}

	d.Repository.All(&models)

	c.JSON(http.StatusOK, models)
}

func (d Post) Get(c *gin.Context) {
	var model = models.Post{}

	model.ID, _ = strconv.Atoi(c.Param("id"))

	d.Repository.Find(&model)
	if model.ID < 1 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "drone not found",
			"status":  http.StatusNotFound,
		})

		return
	}

	c.JSON(http.StatusOK, model)
}
