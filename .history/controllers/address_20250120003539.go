package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/red-star25/advance-go/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

		addresses := make([]models.Address, 0)

		user_id, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			ctx.IndentedJSON(500, "internal server error")
		}

		var c, cancel = context.WithTimeout(context.Background(), time.Second*100)
		defer cancel()

		filter := bson.D{primitive.E{Key: "_id", Value: user_id}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address", Value: addresses}}}}

		_, err = UserCollection.UpdateOne(c, filter, update)
		if err != nil {
			ctx.IndentedJSON(404, "wrong command")
			return
		}
		ctx.Done()
		ctx.IndentedJSON(200, "Successfully Deleted")
	}
}
