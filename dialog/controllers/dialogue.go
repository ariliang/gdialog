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

	// Dialogue history map, {"gastro": <HistoryList>, ...}
	DialogueUtterences []DialogueUtterence

	// Full dialogue history, including multiple model
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
	history, ans := GenDialog(global.Config.DialogCore.Host, history, r.Question, r.DialogModel)
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

// POST /dialoguewx/gastro
// POST /dialoguewx/diabetes
// POST /dialoguewx/mixed-model
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

	// no question just return
	if r.DialogModel != "" && !utils.In([]any{"gastro", "diabetes", "mixed"}, r.DialogModel) {
		return c.JSON(http.StatusOK, utils.Error("All Dialogue Model: gastro, diabetes, mixed"))
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
	history[r.DialogModel] = history[r.DialogModel][utils.Max(len(history[r.DialogModel])-8, 0):] // last 4 round dialogue
	history[r.DialogModel] = append(history[r.DialogModel], DialogueUtterence{"pat", r.Question}) // append "pat:"+question to history

	ans := ""

	history, ans = GenDialog(global.Config.DialogCore.Host, history, r.Question, r.DialogModel)

	log.Println("history: ", history, "ans: ", ans)

	// save session
	hist_byte, _ := json.Marshal(history)
	utils.SetSession(c, data_sess, map[string]any{
		"history": hist_byte,
	})

	if r.DialogModel == "mixed" {
		resp_ans_ls := RespAnsList{RespAns{0, "肠胃病是俗称，为常见病多发病，总发病率约占人口的20%左右。年龄越大，发病率越高，特别是50岁以上的中老年人更为多见。"}, RespAns{1, "糖尿病是一组以高血糖为特征的代谢性疾病。高血糖则是由于胰岛素分泌缺陷或其生物作用受损，或两者兼有引起。长期存在的高血糖，导致各种组织，特别是眼、肾、心脏、血管、神经的慢性损害、功能障碍。"}}
		return c.JSON(http.StatusOK, utils.Success(map[string]any{
			"answer_list": resp_ans_ls,
			"answer_id":   utils.GenerateMD5(fmt.Sprintf("%#v", history)),
		}))
	}

	if r.DialogModel == "gastro" {
		ans = "肠胃病是俗称，为常见病多发病，总发病率约占人口的20%左右。年龄越大，发病率越高，特别是50岁以上的中老年人更为多见。"
	} else if r.DialogModel == "diabetes" {
		ans = "糖尿病是一组以高血糖为特征的代谢性疾病。高血糖则是由于胰岛素分泌缺陷或其生物作用受损，或两者兼有引起。长期存在的高血糖，导致各种组织，特别是眼、肾、心脏、血管、神经的慢性损害、功能障碍。"
	}

	resp_ans_ls := RespAnsList{RespAns{0, ans}}
	return c.JSON(http.StatusOK, utils.Success(map[string]any{
		"answer_list": resp_ans_ls,
		"answer_id":   utils.GenerateMD5(fmt.Sprintf("%#v", history)),
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
