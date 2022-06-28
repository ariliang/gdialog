package dialog

import (
	"gdialog/dialog/controllers"
	"gdialog/dialog/models"
	"gdialog/global"

	"github.com/labstack/echo/v4"
)

// routes in dialog
func Register(e *echo.Echo) {
	// e.GET("/logout", controllers.Logout)
	e.GET("/logoutwx", controllers.LogoutWX)

	// e.POST("/login", controllers.Login)
	e.POST("/loginwx", controllers.LoginWX)
	// e.POST("/register", controllers.Register)
	// e.POST("/dialogue", controllers.Dialogue)
	e.POST("/dialoguewx", controllers.DialogueWX)
	e.POST("/dialoguewx/choose", controllers.DialogueWXChoose)
}

func init() {
	// migrate database
	global.DB.AutoMigrate(&models.User{})
	global.DB.AutoMigrate(&models.UserWX{})
	global.DB.AutoMigrate(&models.Disease{})
}
