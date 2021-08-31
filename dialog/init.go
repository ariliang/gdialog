package dialog

import (
	"gdialog/dialog/models"
	"gdialog/dialog/views"
	"gdialog/global"

	"github.com/labstack/echo/v4"
)

func init() {
	global.DB.AutoMigrate(&models.User{})
	global.DB.AutoMigrate(&models.Disease{})
}

// routes in dialog
func Register(e *echo.Echo) {
	e.GET("/logout", views.Logout)

	e.POST("/login", views.Login)
	e.POST("/register", views.Register)
	e.POST("/dialogue", views.Dialogue)
}
