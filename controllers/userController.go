package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hienphan0111/movie-review-api/database"
	helper "github.com/hienphan0111/movie-review-api/generate-token"
	"github.com/hienphan0111/movie-review-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

func MaskPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var user models.User
		defer cancel()

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Status":  http.StatusBadRequest,
				"Message": "error",
				"Data":    map[string]interface{}{"data": err.Error()},
			})
			return
		}

		regexMatch := bson.M{"$regex": primitive.Regex{Pattern: *user.Username, Options: "i"}}
		emailCount, emailErr := userCollection.CountDocuments(ctx, bson.M{"email": regexMatch})
		usernameMatch := bson.M{"$regex": primitive.Regex{Pattern: *user.Username, Options: "i"}}
		usernameCount, usernameErr := userCollection.CountDocuments(ctx, bson.M{"username": usernameMatch})
		defer cancel()
		if emailErr != nil {
			log.Panic(emailErr)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "error occured while checking for email",
			})
		}
		if usernameErr != nil {
			log.Panic(usernameErr)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "error occured while checking for username",
			})
		}
		if emailCount > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "email already exists",
			})
			return
		}
		if usernameCount > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "username already exists",
			})
			return
		}

		//hash password
		password := MaskPassword(*user.Password)
		user.Password = &password
		user.Created_at = &[]time.Time{time.Now().Local()}[0]
		user.Updated_at = &[]time.Time{time.Now().Local()}[0]
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()
		token, refreshToken, _ := helper.GenerateAllTokens(
			*user.Email,
			*user.Name,
			*user.Username,
			*&user.User_id,
			*user.User_type,
		)
		user.Token = &token
		user.Refresh_token = &refreshToken

		if validationError := validate.Struct(&user); validationError != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Status":  http.StatusBadRequest,
				"Message": "error",
				"Data":    map[string]interface{}{"data": validationError.Error()},
			})
			return
		}

		if validationError := validate.Struct(&user); validationError != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Status":  http.StatusBadRequest,
				"Message": "error",
				"Data":    map[string]interface{}{"data": validationError.Error()},
			})
			return
		}

		//Add new user to the database
		newUser := models.User{
			ID:            user.ID,
			User_id:       user.User_id,
			Name:          user.Name,
			Username:      user.Username,
			Password:      user.Password,
			Email:         user.Email,
			Token:         user.Token,
			Refresh_token: user.Refresh_token,
			User_type:     user.User_type,
			Created_at:    user.Created_at,
			Updated_at:    user.Updated_at,
		}

		result, err := userCollection.InsertOne(ctx, newUser)

		// error message
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Status":  http.StatusInternalServerError,
				"Message": "error",
				"Data":    map[string]interface{}{"data": err.Error()},
			})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Status":  http.StatusInternalServerError,
				"Message": "error",
				"Data":    map[string]interface{}{"data": err.Error()},
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"Status":  http.StatusOK,
			"Message": "success",
			"Data":    map[string]interface{}{"data": result},
		})
	}
}
