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
	}

	RespDialogueCore struct {
		Status string            `json:"status"`
		Result map[string]string `json:"result"`
	}
)

// Given dialogue history list
// Return generated dialogue
func GenDialog(dialog_core_host string, history DialogueHistory, dialog_model string) (DialogueHistory, map[string]string) {

	/*
		dialogue core post data
		{
			"dialog_model": "gastro",
			"history": {
				"gastro": [{"type": "pat", "content", ""}, {"type": "doc", "content": ""}]
			}
		}

		dialogue core return data
		{
			"status": "success",
			"result": {
				"gastro": "",
				"diabetes": ""
			}
		}
	*/

	// prepare dialogue core request data
	req_dialog := ReqDialogueCore{
		dialog_model,
		history,
	}

	resp_dialog := RespDialogueCore{}

	// do request to dialogue core
	req_buf, _ := json.Marshal(req_dialog)
	resp_data, _ := http.Post(dialog_core_host, "application/json", bytes.NewReader(req_buf))

	// process dialogue core response json
	resp_buf, _ := ioutil.ReadAll(resp_data.Body)
	json.Unmarshal(resp_buf, &resp_dialog)

	if resp_dialog.Status != "success" || resp_dialog.Result == nil {
		return history, map[string]string{}
	}

	for dm, ans := range resp_dialog.Result {
		history[dm] = append(history[dm], DialogueUtterence{"doc", ans})
	}

	return history, resp_dialog.Result
}

// simulate wx auth
func SimulateWXAuth(code any) []byte {
	code_list := []any{"123", "456", "789"}
	res := map[string]any{}
	if utils.In(code, code_list) {
		res["session_key"] = "session_key_fjdks"
		res["openid"] = "openid_fkdsfkjd" + code.(string)
	} else {
		res["errcode"] = 40029
		res["errmsg"] = "error rid=fjkdjkfdjkfd"
	}
	res_byte, _ := json.Marshal(res)
	return res_byte
}
