package views

import (
	"encoding/json"
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
	if sess, _ := utils.GetSession(c, "session"); sess.Values["username"] == nil {
		// login session, 7d max age
		sess, _ = utils.Session(c, "session", "/", 3600*24*7)
		utils.SetSession(c, sess, map[string]interface{}{
			"username":  u.Username,
			"logged_in": true,
		})
	} else if sess.Values["logged_in"] == false {
		sess.Values["logged_in"] = true
		sess.Save(c.Request(), c.Response())
	}
	// history session, 20min max age
	if data_sess, err := utils.GetSession(c, "data"); err != nil || data_sess.Values["history"] == nil {
		data_sess, _ = utils.Session(c, "data", "/", 60*20)
		// serialize dict list
		byte_data, _ := json.Marshal([]map[string]string{})
		utils.SetSession(c, data_sess, map[string]interface{}{
			"history": byte_data,
		})
	}
	return c.JSON(http.StatusOK, utils.Success(map[string]interface{}{
		"username": u.Username,
	}))
}
