package handlers

import (
	"testing-api-gateway/api_test/storage"
	"testing-api-gateway/api_test/storage/kv"
	"testing-api-gateway/pkg/etc"
	"encoding/json"
	"log"
	"net/http"
	"net/smtp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/k0kubun/pp"
)

// registrate a user
func RegisterUser(c *gin.Context) {
	var user storage.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user.Id = uuid.NewString()
	user.Email = strings.ToLower(user.Email)
	err := user.Validate()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	userJson, err := json.Marshal(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := kv.Set(user.Id, string(userJson), 1000); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	code := etc.GenerateCode(6)

	auth := smtp.PlainAuth("", "uzbekcodewizard@gmail.com", "bjwz qfzg vxuz xcyr", "smtp.gmail.com")
	err = smtp.SendMail("smtp.gmail.com:587", auth, "uzbekcodewizard", []string{user.Email}, []byte(code))
	if err != nil {
		log.Println(err)
	}

	log.Println("Email sent successfully")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OTP sent to your email",
	})
}

// verify a user
func Verify(c *gin.Context) {
	userCode := c.Param("code")

	if userCode != "12345" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Incorrect code",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
	})
}

// create a user
func CreateUser(c *gin.Context) {
	var user storage.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user.Id = uuid.NewString()

	userJson, err := json.Marshal(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := kv.Set(user.Id, string(userJson), 1000); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

// get a user
func GetUser(c *gin.Context) {
	userID := c.Query("id")
	userString, err := kv.Get(userID)
	pp.Println(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var resp storage.User
	if err := json.Unmarshal([]byte(userString), &resp); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// delete a user
func DeleteUser(c *gin.Context) {
	userId := c.Query("id")
	if err := kv.Delete(userId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user was deleted successfully",
	})
}

// get a list of users
func ListUsers(c *gin.Context) {
	usersStrings, err := kv.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	pp.Println(usersStrings)

	var users []*storage.User
	for _, userString := range usersStrings {
		var user storage.User
		if err := json.Unmarshal([]byte(userString), &user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		users = append(users, &user)
	}

	c.JSON(http.StatusOK, users)
}
