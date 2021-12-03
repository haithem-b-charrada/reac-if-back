package controllers

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"inventory/models"
	"inventory/repositories"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	//"html/template"

	//"log"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Repository *repositories.UserRepository
}

func (u User) Create(c *gin.Context) {
	var model = models.User{}

	c.BindJSON(&model)

	var regex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !regex.MatchString(model.Email) {
		c.JSON(http.StatusBadRequest, "Email not valid")

		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(model.Password), bcrypt.MinCost)

	model.Password = string(hash)

	u.Repository.Save(&model)

	model.Password = ""

	c.JSON(http.StatusOK, model)
	to := []string{model.Email}
	title := "Account verification\n"
	subject := "To confirm your email please pray"
	CreateMail(to, subject, title)
}

func (u User) Update(c *gin.Context) {
	var model = models.User{}
	var request = models.User{}

	model.ID, _ = strconv.Atoi(c.Param("id"))

	u.Repository.Find(&model)
	if model.ID < 1 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "user not found",
			"status":  http.StatusNotFound,
		})

		return
	}

	c.BindJSON(&request)
	model.Email = request.Email

	u.Repository.Save(&model)

	c.JSON(http.StatusOK, model)
}

func (u User) Delete(c *gin.Context) {
	var model = models.User{}

	model.ID, _ = strconv.Atoi(c.Param("id"))

	u.Repository.Find(&model)
	if model.ID < 1 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "user not found",
			"status":  http.StatusNotFound,
		})

		return
	}

	u.Repository.Storage.Delete(&model)

	c.JSON(http.StatusNoContent, "")
}

func (u User) GetAll(c *gin.Context) {
	var models = []models.User{}

	u.Repository.All(&models)

	for i, m := range models {
		m.Password = ""
		models[i] = m
	}

	c.JSON(http.StatusOK, models)
}

func (u User) Get(c *gin.Context) {
	var model = models.User{}

	model.ID, _ = strconv.Atoi(c.Param("id"))

	u.Repository.Find(&model)
	if model.ID < 1 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "user not found",
			"status":  http.StatusNotFound,
		})

		return
	}

	model.Password = ""

	c.JSON(http.StatusOK, model)
}

func (u User) ChangePassword(c *gin.Context) {
	var request map[string]string
	c.BindJSON(&request)

	model := models.User{}
	model.Email = c.Request.Header.Get("User-Email")

	u.Repository.FindByEmail(&model)
	if model.ID < 1 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "user not found",
			"status":  http.StatusNotFound,
		})
	}

	password := model.Password

	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(request["currentPassword"])); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "password incorrect ",
			"status":  http.StatusBadRequest,
		})
		return
	}
	newPassword, _ := bcrypt.GenerateFromPassword([]byte(request["newPassword"]), bcrypt.MinCost)
	if len(request["newPassword"]) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "password should be longer than 8",
		})
		return
	}
	model.Password = string(newPassword)
	u.Repository.Storage.Model(&model).Update("password", string(newPassword))
	c.JSON(http.StatusOK, gin.H{
		"message": "password changed succesfully",
		"status":  http.StatusOK,
	})

}

func tokenGenerator() string {
	b := make([]byte, 10)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
func (u User) ResetPasswordRequest(c *gin.Context) {
	var model = models.User{}
	request := models.User{}
	err := c.BindJSON(&request)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "request not valid",
			"status":  http.StatusBadRequest,
		})

		return
	}
	model.Email = request.Email
	fmt.Println("printing user email", model.Email)
	u.Repository.FindByEmail(&model)
	if model.ID < 1 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "user not found",
			"status":  http.StatusNotFound,
		})

		return
	}
	model.ResetPasswordToken = string(tokenGenerator())
	model.ResetPasswordExpiredDate = time.Now().Add(time.Hour * 2)
	title := "Reset Password\n"
	subject := "Change your passwordlink : url/resetPassword/" + model.ResetPasswordToken
	to := []string{model.Email}
	CreateMail(to, subject, title)
	c.JSON(http.StatusAccepted, gin.H{
		"message": "reset password link sent succesfully",
		"status":  http.StatusAccepted,
	})
	fmt.Println("before staring")
	s := handlejson()
	fmt.Println("after staring", s)
}

//func to handlejson
func handlejson() string {
	absPath, _ := filepath.Abs("../../gateway/config/whitelist.json")
	jsonFile, err := ioutil.ReadFile(absPath)
	if err != nil {
		return err.Error()
	}
	var obj string
	err = json.Unmarshal(jsonFile, &obj)
	return (strings.Replace(obj, "$", "$hhh", -1))
}

//end handing json

//func edit HTML email template

func edit() {
	type Page struct {
		Title string
		Body  []byte
	}

}

//end  edit HTML email template
func (u User) SetNewPassword(c *gin.Context) {
	//Notes to work on :we can take the email from the token
	//Notes to work on :and we will check if we git the email as a variable

	//for testing needs the json body will requires to send the email of the uses to change it password
	//in the  production phase we will disemploy that requirement
	//delete the lines that have a comment after
	request := models.User{}
	model := models.User{}
	err := c.BindJSON(&request)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "request not valid",
			"status":  http.StatusBadRequest,
		})

		return
	}
	model.Email = request.Email      //for test phase
	u.Repository.FindByEmail(&model) //for test phase
	if model.ID < 1 {                //for test phase
		c.JSON(http.StatusNotFound, gin.H{ //for test phase
			"message": "user not found",
			"status":  http.StatusNotFound, //for test phase
		})

		return //for test phase
	} //for test phase

	if len(request.Password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "password should be longer than 8",
			"status":  http.StatusBadRequest,
		})
		fmt.Println("the new password has been detected ")
		return
	}
	newPassword, _ := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.MinCost)
	model.Password = string(newPassword)
	u.Repository.Storage.Model(&model).Update("password", string(newPassword))
	c.JSON(http.StatusOK, gin.H{
		"message": "password changed succesfully",
		"status":  http.StatusOK,
	})

}

func CreateMail(to []string, Subject string, title string) {
	from := "ead.plateform@gmail.com"
	password := "EAD@1234"
	host := "smtp.gmail.com"
	port := "587"
	address := host + ":" + port
	subject := title
	body := Subject
	message := []byte(subject + body)
	auth := smtp.PlainAuth("", from, password, host)
	err := smtp.SendMail(address, auth, from, to, message)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println("Please check your email")
}
