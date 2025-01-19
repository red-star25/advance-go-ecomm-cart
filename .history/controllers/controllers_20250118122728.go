package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/red-star25/advance-go/models"
)

// import "github.com/gin-gonic/gin"

// func HashPassword(password string) string {}

// func VerifyPassword(userPassword string, givenPassword string) (bool, string) {}

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if validationErr := Validate.Struct(user); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr})
		}

		count, err := UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		}

	}
}

// func Login() gin.HandlerFunc {}

// func ProductViewerAdmin() gin.HandlerFunc {}

// func SearchProduct() gin.HandlerFunc {}

// func SearchProductByQuery() gin.HandlerFunc {}
