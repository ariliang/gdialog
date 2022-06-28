package controllers

import (
	"encoding/json"
	"fmt"
	"gdialog/global"
	"gdialog/utils"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	// request json data
	ReqQuestion struct {
		Question    string `json:"question"`
		DialogModel string `json:"dialog_model"`
	}

	// Dialogue utterence struct
	DialogueUtterence struct {
		Type    string `json:"type"`
		Content string `json:"content"`
	}

	// Dialogue history list. [{"type": "pat", "content": ""}, {"type": "doc", "content": ""}]
	DialogueUtterences []DialogueUtterence

	// Full dialogue history, including multiple model. {"gastro": DialogueUtterences, ...}
	DialogueHistory map[string]DialogueUtterences

	RespAns struct {
		No     int    `json:"no"`
		Answer string `json:"answer"`
	}

	RespAnsList []RespAns
)

// POST /dialogue
func Dialogue(c echo.Context) error {
	// get status session whether user logged in
	if sess, err := utils.GetSession(c, "session"); err != nil || sess.Values["username"] == nil || sess.Values["logged_in"] == false {
		return c.JSON(http.StatusUnauthorized, utils.Error("Did not log in"))
	}

	// bind ReqQuestion
	r := ReqQuestion{}
	if err := c.Bind(&r); err != nil {
		return c.JSON(http.StatusOK, utils.Error("Request data error"))
	}

	log.Println("/dialogue request data:", r)

	// no question just return
	if r.Question == "" {
		return c.JSON(http.StatusOK, utils.Error("Input quetion"))
	}

	// get data session
	data_sess, err := utils.GetSession(c, "data")

	if err != nil {
		log.Println("get data session failed...")
		return c.JSON(http.StatusOK, utils.Error("Get data session failed"))
	}

	// read history from session
	history := DialogueHistory{}
	json.Unmarshal(data_sess.Values["history"].([]byte), &history)

	// generate dialog
	history, ans := GenDialog(global.Config.DialogCore.Host, history, r.DialogModel)
	history[r.DialogModel] = history[r.DialogModel][utils.Max(len(history)-8, 0):]                // last 4 round dialogue
	history[r.DialogModel] = append(history[r.DialogModel], DialogueUtterence{"pat", r.Question}) // append "pat:"+question to history

	// save session
	hist_byte, _ := json.Marshal(history)
	utils.SetSession(c, data_sess, map[string]any{
		"history": hist_byte,
	})
	return c.JSON(http.StatusOK, utils.Success(map[string]any{
		"answer": ans,
	}))
}

type (
	ReqDialogueChoose struct {
		Which    int    `json:"which"`
		AnswerID string `json:"answer_id"`
	}
)

func DialogueWXChoose(c echo.Context) error {
	// get status session whether user logged in
	if sess, err := utils.GetSession(c, "session"); err != nil || sess.Values["openid"] == nil || sess.Values["logged_in"] == false {
		return c.JSON(http.StatusUnauthorized, utils.Error("Did not log in"))
	}

	// bind ReqJson
	r := ReqDialogueChoose{-1, ""}
	if err := c.Bind(&r); err != nil {
		return c.JSON(http.StatusOK, utils.Error("Request data error"))
	}

	log.Println("/dialoguewx/choose request data:", r)

	if r.Which == -1 || r.AnswerID == "" {

		return c.JSON(http.StatusOK, utils.Error(fmt.Sprintf("request parameter error: which: %d, answer_id: %s", r.Which, r.AnswerID)))
	}

	return c.JSON(http.StatusOK, utils.Success(map[string]any{}))
}
