package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/red-star25/advance-go/controllers"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("users/signup", controllers.SignUp())
	incomingRoutes.POST("users/login", controllers.Login())
	incomingRoutes.POST("/admin/addProduct", controllers.ProductViewerAdmin())
	incomingRoutes.GET("/users/productView", controllers.SearchProduct())
	incomingRoutes.GET("/users/search", controllers.SearchProductByQuery())
}
