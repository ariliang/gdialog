package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Given dialogue history list
// Return generated dialogue
func GenDialog(history []map[string]string) ([]map[string]string, string) {
	data := map[string]interface{}{
		"history": history,
	}
	byte_data, _ := json.Marshal(data)
	dialogue_server_url := "http://127.0.0.1:5001/"
	resp, _ := http.Post(dialogue_server_url, "application/json", bytes.NewReader(byte_data))
	byte_ans, _ := ioutil.ReadAll(resp.Body)
	resp_json := map[string]string{}
	json.Unmarshal(byte_ans, &resp_json)
	if resp_json["status"] == "success" {
		ans := resp_json["result"]
		history = append(history, map[string]string{"type": "doc", "content": ans})
		return history, ans
	}
	return history, "请重新输入"
}

// simulate wx auth
func SimulateWXAuth(code interface{}) []byte {
	code_list := []interface{}{"123", "456", "789"}
	res := map[string]interface{}{}
	if In(code_list, code) {
		res["session_key"] = "session_key_fjdks"
		res["openid"] = "openid_fkdsfkjd" + code.(string)
	} else {
		res["errcode"] = 40029
		res["errmsg"] = "error rid=fjkdjkfdjkfd"
	}
	res_byte, _ := json.Marshal(res)
	return res_byte
}
