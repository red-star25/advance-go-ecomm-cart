package controllers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/red-star25/advance-go/database"
	"github.com/red-star25/advance-go/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	prodCollection *mongo.Collection
	userCollection *mongo.Collection
}

func NewApplication(prodCollection, userCollection *mongo.Collection) *Application {
	return &Application{
		prodCollection: prodCollection,
		userCollection: userCollection,
	}
}

func (app *Application) AddToCart() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		productQueryID := ctx.Query("id")
		if productQueryID == "" {
			log.Println("product id is empty")

			_ = ctx.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}

		userQueryID := ctx.Query("userID")
		if userQueryID == "" {
			log.Println("user id is empty")

			_ = ctx.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}

		productID, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			log.Println(err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var c, cancel = context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		if err := database.AddProductToCart(c, app.prodCollection, app.userCollection, productID, userQueryID); err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, err)
			return

		}
		ctx.IndentedJSON(200, "Successfully added to cart")

	}
}

func (app *Application) RemoveItem() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		productQueryID := ctx.Query("id")
		if productQueryID == "" {
			log.Println("product id is empty")

			_ = ctx.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}

		userQueryID := ctx.Query("userID")
		if userQueryID == "" {
			log.Println("user id is empty")

			_ = ctx.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}

		productID, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			log.Println(err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var c, cancel = context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		if err := database.RemoveCartItem(c, app.prodCollection, app.userCollection, productID, userQueryID); err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		ctx.IndentedJSON(200, "Successfully removed item from cart")
	}
}

func GetItemFromCart() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.Query("id")

		if userID == "" {
			ctx.Header("Content-Type", "application/json")
			ctx.JSON(http.StatusNotFound, gin.H{"error": "invalid id"})
			ctx.Abort()
			return
		}

		user_id, _ := primitive.ObjectIDFromHex(userID)
		var c, cancel = context.WithTimeout(context.Background(), time.Second*100)
		defer cancel()

		var filledCart models.User
		err := UserCollection.FindOne(c, bson.D{primitive.E{Key: "_id", Value: user_id}}).Decode(&filledCart)
		if err != nil {
			log.Println(err)
			ctx.IndentedJSON(500, "not found")
			return
		}

		filterMatch := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: userID}}}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$usercart"}}}}
		grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$usercart.price"}}}}}}

		pointCursor, err := UserCollection.Aggregate(c, mongo.Pipeline{filterMatch, unwind, grouping})
		if err != nil {
			fmt.Println(err)
			return
		}
		var listing []bson.M
		if err := pointCursor.All(c, &listing); err != nil {
			log.Println(err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}

		for _, json := range listing {
			ctx.IndentedJSON(200, json["total"])
			ctx.IndentedJSON(200, filledCart.UserCart)
		}

		ctx.Done()

	}
}

func (app *Application) BuyFromCart() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userQueryID := ctx.Query("id")

		if userQueryID == "" {
			log.Panic("user id is empty")
			ctx.AbortWithError(http.StatusBadRequest, errors.New("userID is empty"))
		}

		c, cancel := context.WithTimeout(context.Background(), time.Second*100)
		defer cancel()

		if err := database.BuyItemFromCart(c, app.userCollection, userQueryID); err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		ctx.IndentedJSON(200, "Successfully placed the order")
	}
}

func (app *Application) InstantBuy() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		productQueryID := ctx.Query("id")
		if productQueryID == "" {
			log.Println("product id is empty")

			_ = ctx.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}

		userQueryID := ctx.Query("userID")
		if userQueryID == "" {
			log.Println("user id is empty")

			_ = ctx.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}

		productID, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			log.Println(err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var c, cancel = context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		if err := database.InstantBuyer(c, app.prodCollection, app.userCollection, productID, userQueryID); err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		ctx.IndentedJSON(200, "Successfully placed the order")
	}
}
