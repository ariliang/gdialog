package controllers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"gdialog/utils"
)

type (
	ReqDialogueCore struct {
		DialogModel string `json:"dialog_model"`
		History     any    `json:"history"`
		Question    string `json:"question"`
	}

	RespDialogueCore map[string]string
)

// Given dialogue history list
// Return generated dialogue
func GenDialog(dialog_core_host string, history DialogueHistory, question string, dialog_model string) (DialogueHistory, string) {

	/*
		dialogue core post data
		{
			"dialog_model": "gastro",
			"history": {
				"gastro": [{"type": "pat", "content", ""}, {"type": "doc", "content": ""}]
			},
			"question": ""
		}

		dialogue core return data
		{
			"status": "success",
			"result": "ans"
		}
	*/

	// prepare dialogue core request data
	req_dialog := ReqDialogueCore{
		dialog_model,
		history,
		question,
	}

	resp_dialog := RespDialogueCore{}

	// do request to dialogue core
	req_buf, _ := json.Marshal(req_dialog)
	resp_data, _ := http.Post(dialog_core_host, "application/json", bytes.NewReader(req_buf))

	// process dialogue core response json
	resp_buf, _ := ioutil.ReadAll(resp_data.Body)
	json.Unmarshal(resp_buf, &resp_dialog)

	if resp_dialog["status"] != "success" || resp_dialog["result"] == "" {
		return history, ""
	}

	ans := resp_dialog["result"]
	history[dialog_model] = append(history[dialog_model], DialogueUtterence{"pat", question})
	history[dialog_model] = append(history[dialog_model], DialogueUtterence{"doc", ans})

	return history, ans
}

// simulate wx auth
func SimulateWXAuth(code any) []byte {
	code_list := []any{"123", "456", "789"}
	res := map[string]any{}
	if utils.In(code_list, code) {
		res["session_key"] = "session_key_fjdks"
		res["openid"] = "openid_fkdsfkjd" + code.(string)
	} else {
		res["errcode"] = 40029
		res["errmsg"] = "error rid=fjkdjkfdjkfd"
	}
	res_byte, _ := json.Marshal(res)
	return res_byte
}
