package main

import (
	"fmt"
	"gdialog/dialog"
	"gdialog/global"

	"github.com/BurntSushi/toml"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func StartServer() {
	// read server config
	server_conf := new(global.ServerConfig)
	toml.DecodeFile("conf/server.toml", &server_conf)
	// create server
	e := echo.New()
	// set session
	e.Use(middleware.Recover())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("ibiauhsihsow"))))
	// register routes
	dialog.Register(e)
	// start server
	e.Logger.Info(e.Start(fmt.Sprintf("%s:%d", server_conf.Host, server_conf.Port)))
}

func main() {
	StartServer()
	// Test()
}
