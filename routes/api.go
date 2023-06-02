package routes

import (
	historiesController "sahamrakyat_test/histories/controller"
	ordersController "sahamrakyat_test/orders/controller"
	usersController "sahamrakyat_test/users/controller"

	"github.com/labstack/echo/v4"
)

func Init(app *echo.Echo) {
	apiGroup := app.Group("/api")
	// v1
	apiv1Group := apiGroup.Group("/v1")
	ordersGroup := apiv1Group.Group("/orders")
	ordersGroup.GET("/", ordersController.GetAll)
	ordersGroup.GET("/:id", ordersController.Get)
	ordersGroup.POST("/", ordersController.Create)
	ordersGroup.PUT("/:id", ordersController.Update)
	ordersGroup.DELETE("/:id", ordersController.Delete)
	usersGroup := apiv1Group.Group("/users")
	usersGroup.GET("/", usersController.GetAll)
	usersGroup.GET("/:id", usersController.Get)
	usersGroup.POST("/", usersController.Create)
	usersGroup.PUT("/:id", usersController.Update)
	usersGroup.DELETE("/:id", usersController.Delete)
	orderHistoriesGroup := apiv1Group.Group("/histories")
	orderHistoriesGroup.GET("/", historiesController.GetAll)
	orderHistoriesGroup.GET("/:id", historiesController.Get)
	orderHistoriesGroup.POST("/", historiesController.Create)
	orderHistoriesGroup.PUT("/:id", historiesController.Update)
	orderHistoriesGroup.DELETE("/:id", historiesController.Delete)
	// v2
	// Note: If you had breaking change in API, use new version for preserving old API while creating new one
	// ...
}
