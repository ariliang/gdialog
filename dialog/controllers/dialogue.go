package controllers

import (
	"encoding/json"
	"fmt"
	"gdialog/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	// request json data
	ReqJson struct {
		Question string `json:"question"`
	}

	// history list, list of map
	HistoryList []map[string]string
)

// POST /dialogue
func Dialogue(c echo.Context) error {
	// get status session whether user logged in
	if sess, err := utils.GetSession(c, "session"); err != nil || sess.Values["username"] == nil || sess.Values["logged_in"] == false {
		return c.JSON(http.StatusUnauthorized, utils.Error("Did not log in"))
	}

	// bind ReqJson
	r := ReqJson{}
	if err := c.Bind(&r); err != nil {
		return c.JSON(http.StatusOK, utils.Error("Request data error"))
	}

	// no question just return
	if r.Question == "" {
		return c.JSON(http.StatusOK, utils.Error("Input quetion"))
	}

	history := HistoryList{}
	// get data session
	data_sess, err := utils.GetSession(c, "data")

	if err != nil {
		fmt.Println("failed")
		return c.JSON(http.StatusOK, utils.Error("Get data session failed"))
	}

	// read history from session
	json.Unmarshal(data_sess.Values["history"].([]byte), &history)
	history = history[utils.Max(len(history)-8, 0):]                                   // last 4 round dialogue
	history = append(history, map[string]string{"type": "pat", "content": r.Question}) // append "pat:"+question to history
	// generate dialog
	history, ans := utils.GenDialog(history)
	// save session
	hist_byte, _ := json.Marshal(history)
	utils.SetSession(c, data_sess, map[string]interface{}{
		"history": hist_byte,
	})
	return c.JSON(http.StatusOK, utils.Success(map[string]interface{}{
		"answer": ans,
	}))
}

// POST /dialoguewx
func DialogueWX(c echo.Context) error {
	// get status session whether user logged in
	if sess, err := utils.GetSession(c, "session"); err != nil || sess.Values["openid"] == nil || sess.Values["logged_in"] == false {
		return c.JSON(http.StatusUnauthorized, utils.Error("Did not log in"))
	}

	// bind ReqJson
	r := ReqJson{}
	if err := c.Bind(&r); err != nil {
		return c.JSON(http.StatusOK, utils.Error("Request data error"))
	}

	// no question just return
	if r.Question == "" {
		return c.JSON(http.StatusOK, utils.Error("Input quetion"))
	}

	history := HistoryList{}
	// get data session
	data_sess, err := utils.GetSession(c, "data")

	if err != nil {
		return c.JSON(http.StatusOK, utils.Error("Get data session failed"))
	}

	// read history from session
	json.Unmarshal(data_sess.Values["history"].([]byte), &history)
	history = history[utils.Max(len(history)-8, 0):]                                   // last 4 round dialogue
	history = append(history, map[string]string{"type": "pat", "content": r.Question}) // append "pat:"+question to history
	// generate dialog
	history, ans := utils.GenDialog(history)
	// save session
	hist_byte, _ := json.Marshal(history)
	utils.SetSession(c, data_sess, map[string]interface{}{
		"history": hist_byte,
	})
	return c.JSON(http.StatusOK, utils.Success(map[string]interface{}{
		"answer": ans,
	}))
}