package routes

import (
	"app/app/controller"

	"github.com/gin-gonic/gin"
)

func Prize(router *gin.RouterGroup) {
	// Get the *bun.DB instance from config
	ctl := controller.New() // Pass the *bun.DB to the controller
	// md := middleware.AuthMiddleware()
	// log := middleware.NewLogResponse()
	prize := router.Group("")
	{
		prize.POST("/create", ctl.PrizeCtl.Create)
		prize.PATCH("/:id", ctl.PrizeCtl.Update)
		prize.GET("/list", ctl.PrizeCtl.List)
		prize.GET("/:id", ctl.PrizeCtl.Get)
		prize.DELETE("/:id", ctl.PrizeCtl.Delete)
		prize.POST("/upload", ctl.PrizeCtl.UploadImage)

	}
}
