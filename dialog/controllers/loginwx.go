package controllers

import (
	"encoding/json"
	"fmt"
	"gdialog/dialog/models"
	"gdialog/global"
	"gdialog/utils"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	authUrl = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
)

type (
	// from client
	reqData struct {
		Code     string
		Username string
		Gender   uint8
	}
	// from wx server
	respData struct {
		// failed
		ErrCode int
		ErrMsg  string
		// success
		SessionKey string `json:"session_key"`
		OpenID     string
	}
)

// POST /loginwx
func LoginWX(c echo.Context) error {
	// bind User
	req_data := new(reqData)
	if err := c.Bind(req_data); err != nil {
		return c.JSON(http.StatusUnauthorized, utils.Error("Authorization Failure"))
	}
	// auth from wx server
	auth_req := fmt.Sprintf(authUrl, global.Config.WX.AppId, global.Config.WX.AppIdSecret, req_data.Code)
	resp, _ := http.Get(auth_req)
	var resp_byte []byte
	// debug on to enable simulated auth
	debug := false
	if !debug {
		resp_byte, _ = ioutil.ReadAll(resp.Body)
	} else {
		resp_byte = utils.SimulateWXAuth(req_data.Code)
	}
	// parse json body to struct
	resp_data := respData{}
	json.Unmarshal(resp_byte, &resp_data)
	// auth failed
	if resp.StatusCode != 200 || resp_data.OpenID == "" {
		return c.JSON(http.StatusUnauthorized, utils.Error("Authorization Failure"))
	}
	// bind user
	u := models.UserWX{
		OpenID:   resp_data.OpenID,
		Username: req_data.Username,
		Gender:   req_data.Gender,
	}
	fmt.Println(u)
	// create user if not existed
	if !u.Exists() {
		if err := u.Save(); err != nil {
			return c.JSON(http.StatusOK, utils.Error("Create user failed"))
		}
	}
	// read from session
	if sess, _ := utils.GetSession(c, "session"); sess.Values["openid"] == nil {
		// login session, 7d max age
		sess, _ = utils.Session(c, "session", "/", 3600*24*7)
		utils.SetSession(c, sess, map[string]interface{}{
			"openid":      u.Username,
			"session_key": resp_data.SessionKey,
			"logged_in":   true,
		})
	} else if sess.Values["logged_in"] == false {
		utils.SetSession(c, sess, map[string]interface{}{
			"logged_in": true,
		})
	}
	// history session, 20min max age
	if data_sess, err := utils.GetSession(c, "data"); err != nil || data_sess.Values["history"] == nil {
		data_sess, _ = utils.Session(c, "data", "/", 60*20)
		// serialize dict list
		byte_data, _ := json.Marshal(HistoryList{})
		utils.SetSession(c, data_sess, map[string]interface{}{
			"history": byte_data,
		})
	}
	// login succeed
	return c.JSON(http.StatusOK, utils.Success(map[string]interface{}{
		"username": u.Username,
	}))
}
