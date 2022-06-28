package controllers

import (
	"encoding/json"
	"fmt"
	"gdialog/global"
	"gdialog/utils"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

// POST /dialoguewx
func DialogueWX(c echo.Context) error {
	// get status session whether user logged in
	if sess, err := utils.GetSession(c, "session"); err != nil || sess.Values["openid"] == nil || sess.Values["logged_in"] == false {
		return c.JSON(http.StatusUnauthorized, utils.Error("Did not log in"))
	}

	// bind ReqJson
	r := ReqQuestion{}
	if err := c.Bind(&r); err != nil {
		return c.JSON(http.StatusOK, utils.Error("Request data error"))
	}

	log.Println("/dialoguewx request data:", r)

	// no question just return
	if r.Question == "" {
		return c.JSON(http.StatusOK, utils.Error("Input quetion"))
	}

	dms := global.Config.DialogCore.DialogModel

	// no specified model just return
	if all_dms := append(dms, "mixed"); r.DialogModel != "" && !utils.In(r.DialogModel, utils.StrListToAny(all_dms)) {
		return c.JSON(http.StatusOK, utils.Error(fmt.Sprintf("All Dialogue Model: %s", all_dms)))
	}

	// get data session
	data_sess, err := utils.GetSession(c, "data")

	if err != nil {
		log.Println("get data session failed...")
		return c.JSON(http.StatusOK, utils.Error("Get data session failed, please relog in"))
	}

	// read history from session
	history := DialogueHistory{}

	if r.DialogModel != "mixed" {
		json.Unmarshal(data_sess.Values["history"].([]byte), &history)

		history[r.DialogModel] = history[r.DialogModel][utils.Max(len(history[r.DialogModel])-8, 0):] // last 4 round dialogue
		history[r.DialogModel] = append(history[r.DialogModel], DialogueUtterence{"pat", r.Question}) // append "pat:"+question to history
	} else {
		json.Unmarshal(data_sess.Values["mixed_history"].([]byte), &history)

		for _, dm := range dms {
			history[dm] = history[dm][utils.Max(len(history[dm])-8, 0):]            // last 4 round dialogue
			history[dm] = append(history[dm], DialogueUtterence{"pat", r.Question}) // append "pat:"+question to history
		}
	}

	log.Println("history in data session:", history)

	answers := map[string]string{}

	history, answers = GenDialog(global.Config.DialogCore.Host, history, r.DialogModel)

	log.Println("history: ", history, "answers: ", answers)

	// save session
	hist_byte, _ := json.Marshal(history)

	if r.DialogModel != "mixed" {
		utils.SetSession(c, data_sess, map[string]any{
			"history": hist_byte,
		})

	} else {
		utils.SetSession(c, data_sess, map[string]any{
			"mixed_history": hist_byte,
		})
	}

	answers = lo.PickByKeys(answers, dms)

	resp_ans_ls := RespAnsList{}
	i := 0
	for _, v := range answers {
		resp_ans_ls = append(resp_ans_ls, RespAns{i, v})
		i++
	}

	log.Println("/dialoguewx response data:", resp_ans_ls)

	return c.JSON(http.StatusOK, utils.Success(map[string]any{
		"answer_list": resp_ans_ls,
		"answer_id":   utils.GenerateMD5(fmt.Sprintf("%#v", history)),
	}))
}
