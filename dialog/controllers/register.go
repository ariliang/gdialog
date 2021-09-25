package controllers

import (
	"gdialog/dialog/models"
	"gdialog/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

// POST /register
func Register(c echo.Context) error {
	// bind json
	u := new(models.User)
	if err := c.Bind(u); err != nil {
		return c.JSON(http.StatusOK, utils.Error("Check your info"))
	}
	// user exists
	if u.Exists() {
		return c.JSON(http.StatusOK, utils.Error("User existed"))
	}
	// validate info
	if ok, msg := u.Validate(); !ok {
		return c.JSON(http.StatusOK, utils.Error(msg))
	}
	// validate disease
	d := new(models.Disease)
	d.Disease = u.Disease
	if ok, msg := d.Valid(); !ok {
		return c.JSON(http.StatusOK, utils.Error(msg))
	}
	// save User
	if err := u.Save(); err != nil {
		return c.JSON(http.StatusOK, utils.Error("Register user failed"))
	}
	// save disease
	if d.Disease != "" {
		// get newly created user id
		u.GetUser()
		d.UserID = u.UserID
		// save Disease
		d.Save()
	}
	return c.JSON(http.StatusOK, utils.Success(map[string]interface{}{"username": u.Username}))
}
