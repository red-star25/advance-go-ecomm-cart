package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// import "github.com/gin-gonic/gin"

// func AddAddress() gin.HandlerFunc{}

// func EditAddress() gin.HandlerFunc{}

// func EditWorkAddress() gin.HandlerFunc{}

func DeleteAddress() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.Query("id")

		if userID == "" {
			ctx.Header("Content-Type", "application/json")
			ctx.JSON(http.StatusNotFound, gin.H{"error": "invalid search index"})
			ctx.Abort()
			return
		}
	}
}
