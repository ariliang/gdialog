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
	resp, _ := http.Post("http://127.0.0.1:5001/", "application/json", bytes.NewReader(byte_data))
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
