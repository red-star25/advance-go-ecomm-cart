package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/red-star25/advance-go/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddAddress() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.Query("id")

		if userID == "" {
			ctx.Header("Content-Type", "application/json")
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Invalid code"})
			ctx.Abort()
			return
		}

		address, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			ctx.IndentedJSON(500, "internal server error")
		}

		var addresses models.Address

		addresses.AdressID = primitive.NewObjectID()

		if err := ctx.BindJSON(&addresses); err != nil {
			ctx.IndentedJSON(http.StatusNotAcceptable, err.Error())
		}

		var c, cancel = context.WithTimeout(context.Background(), time.Second*100)
		defer cancel()

		matchFilter := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: address}}}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$address"}}}}
		group := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$address_id"}, {Key: "count", Value: bson.D{primitive.E{Key: "$sum", Value: 1}}}}}}

		pointCursor, err := UserCollection.Aggregate(c, mongo.Pipeline{matchFilter, unwind, group})
		if err != nil {
			ctx.IndentedJSON(500, "internal server error")
		}

		var addressInfo []bson.M
		if err = pointCursor.All(c, &addressInfo); err != nil {
			panic(err)
		}

		var size int32
		for _, address_no := range addressInfo {
			count := address_no["count"]
			size = count.(int32)
		}
		if size < 2 {
			filter := bson.D{primitive.E{Key: "_id", Value: address}}
			update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "address", Value: addresses}}}}
			_, err := UserCollection.UpdateOne(c, filter, update)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			ctx.IndentedJSON(400, "not allowed")
		}
		defer cancel()
		ctx.Done()
	}
}

func EditHomeAddress() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.Query("id")

		if userID == "" {
			ctx.Header("Content-Type", "application/json")
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Invalid"})
			ctx.Abort()
			return
		}

		user_id, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			ctx.IndentedJSON(500, "internal server error")
		}
		var editAddress models.Address
		if err := ctx.BindJSON(&editAddress); err != nil {
			ctx.IndentedJSON(http.StatusBadRequest, err.Error())
		}
		var c, cancel = context.WithTimeout(context.Background(), time.Second*100)
		defer cancel()

		filter := bson.D{primitive.E{Key: "_id", Value: user_id}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address.0.house_name", Value: editAddress.House}, {Key: "address.0.street_name", Value: editAddress.Street}, {Key: "address.0.city_name", Value: editAddress.City}, {Key: "address.0.pin_code", Value: editAddress.Pincode}}}}

		_, err = UserCollection.UpdateOne(c, filter, update)
		if err != nil {
			ctx.IndentedJSON(500, "something went wrong")
			return
		}
		defer cancel()
		ctx.Done()
		ctx.IndentedJSON(200, "Successfully updated Home address")
	}
}

func EditWorkAddress() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.Query("id")

		if userID == "" {
			ctx.Header("Content-Type", "application/json")
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Invalid"})
			ctx.Abort()
			return
		}

		user_id, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			ctx.IndentedJSON(500, "internal server error")
		}
		var editAddress models.Address
		if err := ctx.BindJSON(&editAddress); err != nil {
			ctx.IndentedJSON(http.StatusBadRequest, err.Error())
		}
		var c, cancel = context.WithTimeout(context.Background(), time.Second*100)
		defer cancel()

		filter := bson.D{primitive.E{Key: "_id", Value: user_id}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address.1.house_name", Value: editAddress.House}, {Key: "address.1.street_name", Value: editAddress.Street}, {Key: "address.1.city_name", Value: editAddress.City}, {Key: "address.1.pin_code", Value: editAddress.Pincode}}}}

		_, err = UserCollection.UpdateOne(c, filter, update)
		if err != nil {
			ctx.IndentedJSON(500, "something went wrong")
			return
		}
		defer cancel()
		ctx.Done()
		ctx.IndentedJSON(200, "Successfully updated Home address")

	}
}

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
