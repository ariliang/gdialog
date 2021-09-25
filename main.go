package main

import (
	"fmt"
	"gdialog/dialog"
	"gdialog/global"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func StartServer() {
	// create server
	e := echo.New()
	// set session
	e.Use(middleware.Recover())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(global.Config.Web.SessionSecret))))
	// register apps
	dialog.Register(e)
	// start server
	e.Logger.Info(e.Start(fmt.Sprintf("%s:%d", global.Config.Server.Host, global.Config.Server.Port)))
}

func main() {
	StartServer()
	// Test()
}
