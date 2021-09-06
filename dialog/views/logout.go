package views

import (
	"gdialog/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Logout(c echo.Context) error {
	username := c.QueryParam("username")
	sess, err := utils.GetSession(c, "session")
	if err != nil || sess.Values["username"] != username || sess.Values["logged_in"] == false {
		return c.JSON(http.StatusOK, utils.Error("Did not log in"))
	}
	utils.SetSession(c, sess, map[string]interface{}{
		"logged_in": false,
	})
	return c.JSON(http.StatusOK, utils.Success(nil))
}
