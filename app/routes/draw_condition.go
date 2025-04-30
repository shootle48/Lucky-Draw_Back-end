package routes

import (
	"app/app/controller"

	"github.com/gin-gonic/gin"
)

func DrawCondition(router *gin.RouterGroup) {
	// Get the *bun.DB instance from config
	ctl := controller.New() // Pass the *bun.DB to the controller
	// md := middleware.AuthMiddleware()
	// log := middleware.NewLogResponse()
	draw_condition := router.Group("")
	{
		draw_condition.POST("/create", ctl.DrawConditionCtl.Create)
		draw_condition.PATCH("/:id", ctl.DrawConditionCtl.Update)
		draw_condition.GET("/list", ctl.DrawConditionCtl.List)
		draw_condition.GET("/:id", ctl.DrawConditionCtl.Get)
		draw_condition.DELETE("/:id", ctl.DrawConditionCtl.Delete)
		draw_condition.POST("/preview", ctl.DrawConditionCtl.PreviewPlayer)
		draw_condition.GET("/GetDrawConditionPreview/:id", ctl.DrawConditionCtl.GetDrawConditionPreview)
	}
}
