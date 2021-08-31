package views

import (
	"gdialog/dialog/models"
	"gdialog/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

// POST(json)
func Login(c echo.Context) error {
	// bind User
	u := new(models.User)
	if err := c.Bind(u); err != nil {
		return c.JSON(http.StatusOK, utils.Error("Check your info"))
	}
	// if User not exists
	if !u.Exists() {
		return c.JSON(http.StatusOK, utils.Error("User does not exist"))
	}
	// user and password do not match
	if !u.ValidLogin() {
		return c.JSON(http.StatusOK, utils.Error("Password does not match"))
	}
	// login session, 7d max age
	sess, _ := utils.Session(c, "session", "/", 3600*24*7)
	utils.SetSession(c, sess, map[string]interface{}{
		"username":  u.Username,
		"logged_in": true,
	})
	// history session, 20min max age
	data_sess, _ := utils.Session(c, "data", "/", 60*20)
	utils.SetSession(c, data_sess, map[string]interface{}{
		"history": []string{},
	})
	return c.JSON(http.StatusOK, utils.Success(map[string]interface{}{
		"username": u.Username,
	}))
}