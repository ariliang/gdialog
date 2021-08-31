package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Given dialogue history list
// Return generated dialogue
func GenDialog(history []string) ([]string, string) {
	data := map[string]interface{}{
		"history": history,
	}
	byte_data, _ := json.Marshal(data)
	resp, _ := http.Post("http://127.0.0.1:5001/", "application/json", bytes.NewReader(byte_data))
	byte_ans, _ := ioutil.ReadAll(resp.Body)
	ans := string(byte_ans)
	history = append(history, fmt.Sprintf("doc:%s", ans))
	return history, ans
}
