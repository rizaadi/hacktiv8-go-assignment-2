package routers

import (
	"hacktiv8-go-assignment-2/controllers"

	"github.com/gin-gonic/gin"
)

func StartServer() (router *gin.Engine) {
	router = gin.Default()

	router.POST("/orders", controllers.CreateOrder)
	router.GET("/orders", controllers.GetOrders)
	router.GET("/orders/:id", controllers.GetOrderById)
	router.PUT("/orders/:id", controllers.UpdateOrderById)
	router.DELETE("/orders/:id", controllers.DeleteOrderById)
	return
}
