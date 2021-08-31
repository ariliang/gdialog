package main

import (
	"fmt"
	"gdialog/global"
	"gdialog/utils"
	"strings"
)

func Test() {
	SliceTest()
}

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
	history, _ := map[string]interface{}{}["history"].([]string)
	history = history[utils.Max(len(history)-8, 0):]
	history = append(history, "hello")
	fmt.Println(history)
}
