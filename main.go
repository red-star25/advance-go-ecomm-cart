package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/red-star25/advance-go/controllers"
	"github.com/red-star25/advance-go/database"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	app := controllers.NewApplication(database.ProductData(database.Client, "Products"), database.UserData(database.Client, "Users"))

	router := gin.New()
	router.Use(gin.Logger())

	router.UserRoutes(router)
	router.Use(middleware.Authentication())

	router.GET("/addToCart", app.AddToCart())
	router.GET("/removeItem", app.RemoveItem())
	router.GET("/cartCheckout", app.CartCheckout())
	router.GET("/instantBuy", app.InstantBuy())

	log.Fatal(router.Run(":" + port))
}
