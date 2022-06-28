package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gdialog/global"
	"gdialog/utils"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strings"
)

func DBTest() {
	fmt.Println("hello")
	db := global.DB
	fmt.Println(db)
}

func SliceTest() {
	var s []string
	fmt.Println(s)

	ds := strings.Split(";糖尿病", ";")
	fmt.Println(ds)

	for _, v := range ds {
		fmt.Println(v)
	}

	s = append(s, "hello")
	fmt.Println(s)

	fmt.Println(strings.Join(s, ";"))
	fmt.Println("============================")
	history, _ := map[string]any{}["history"].([]string)
	history = history[utils.Max(len(history)-8, 0):]
	history = append(history, "hello")
	fmt.Println(history)
}

func RequestTest() {
	data := map[string]any{
		"history": []string{"hello"},
	}
	byte_data, _ := json.Marshal(data)
	resp, _ := http.Post("http://127.0.0.1:5001/", "application/json", bytes.NewReader(byte_data))
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserProfile struct {
	// inherit User
	User
	Age    int    `json:"age"`
	Gender string `json:"gender"`
	Phone  string `json:"phone"`
}

func (u User) Valid() bool {
	if u.Username != "" && u.Password != "" {
		return true
	}
	return false
}

func (u UserProfile) Valid() bool {
	// override
	// call parent's function
	return u.User.Valid() && u.Age >= 0 && u.Age < 200
}

func InheritanceTest() {
	// init inherited fields
	user := UserProfile{
		User: User{
			Username: "fds",
			Password: "f",
		},
		Age: 3,
	}
	fmt.Println(user)

	// created by json
	user = UserProfile{}
	json.Unmarshal([]byte("{\"username\": \"df\",\"password\": \"p\",\"age\": -3,\"gender\": \"F\",\"phone\": \"34343\"}"), &user)

	// compare types
	if reflect.TypeOf(user) == reflect.TypeOf(UserProfile{}) {
		fmt.Println("typeof user is: ", reflect.TypeOf(UserProfile{}))
	}
	fmt.Println(user)
	fmt.Println(user.Username)
	fmt.Println(user.User.Username)
	fmt.Println(user.Gender)
	fmt.Println(user.Valid())

	// get field's tag
	fmt.Println(reflect.TypeOf(user).Name())

}

func Test() {
	m := map[string]string{"a": "aa"}
	log.Println(m)
}
