package routes

import (
	"app/app/controller"

	"github.com/gin-gonic/gin"
)

func Player(router *gin.RouterGroup) {
	// Get the *bun.DB instance from config
	ctl := controller.New() // Pass the *bun.DB to the controller
	// md := middleware.AuthMiddleware()
	// log := middleware.NewLogResponse()
	player := router.Group("")
	{
		player.POST("/create", ctl.PlayerCtl.Create)
		player.PATCH("/:id", ctl.PlayerCtl.Update)
		player.GET("/list", ctl.PlayerCtl.List)
		player.GET("/:id", ctl.PlayerCtl.Get)
		player.DELETE("/:id", ctl.PlayerCtl.Delete)
		player.POST("/import", ctl.PlayerCtl.ImportCSV)
	}
}
