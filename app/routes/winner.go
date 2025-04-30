package routes

import (
	"app/app/controller"

	"github.com/gin-gonic/gin"
)

func WinnerRoutes(router *gin.RouterGroup) {
	// Get the *bun.DB instance from config
	ctl := controller.New() // Pass the *bun.DB to the controller
	// md := middleware.AuthMiddleware()
	// log := middleware.NewLogResponse()
	winner := router.Group("")
	{
		winner.POST("/create", ctl.WinnerCtl.Create)
		winner.PATCH("/:id", ctl.WinnerCtl.Update)
		winner.GET("/list", ctl.WinnerCtl.List)
		winner.GET("/:id", ctl.WinnerCtl.Get)
		winner.DELETE("/:id", ctl.WinnerCtl.Delete)
		winner.GET("/room/:room_id", ctl.WinnerCtl.DashboardByRoomID)

	}
}
