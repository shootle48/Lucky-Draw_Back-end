package routes

import (
	"app/app/controller"

	"github.com/gin-gonic/gin"
)

func Room(router *gin.RouterGroup) {
	// Get the *bun.DB instance from config
	ctl := controller.New() // Pass the *bun.DB to the controller
	// md := middleware.AuthMiddleware()
	// log := middleware.NewLogResponse()
	room := router.Group("")
	{
		room.POST("/create", ctl.RoomCtl.Create)
		room.PATCH("/:id", ctl.RoomCtl.Update)
		room.GET("/list", ctl.RoomCtl.List)
		room.GET("/:id", ctl.RoomCtl.Get)
		room.DELETE("/:id", ctl.RoomCtl.Delete)
		room.GET("/:id/all", ctl.RoomCtl.ListAll)
	}
}
