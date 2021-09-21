package dialog

import (
	"gdialog/dialog/controllers"
	"gdialog/dialog/models"
	"gdialog/global"

	"github.com/labstack/echo/v4"
)

func init() {
	global.DB.AutoMigrate(&models.User{})
	global.DB.AutoMigrate(&models.Disease{})
}

// routes in dialog
func Register(e *echo.Echo) {
	e.GET("/logout", controllers.Logout)

	e.POST("/login", controllers.Login)
	e.POST("/register", controllers.Register)
	e.POST("/dialogue", controllers.Dialogue)
}
