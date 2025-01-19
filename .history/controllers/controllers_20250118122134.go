package controllers

import "github.com/gin-gonic/gin"

// import "github.com/gin-gonic/gin"

// func HashPassword(password string) string {}

// func VerifyPassword(userPassword string, givenPassword string) (bool, string) {}

func Signup() gin.HandlerFunc {
	return func(c *gin.Context){
		var ctx, cancel:= context.WithTimeout(context.Background(), 100 * time.Second)
	}
}

// func Login() gin.HandlerFunc {}

// func ProductViewerAdmin() gin.HandlerFunc {}

// func SearchProduct() gin.HandlerFunc {}

// func SearchProductByQuery() gin.HandlerFunc {}
