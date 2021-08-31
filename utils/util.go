package utils

import (
	"fmt"
)

// Given dialogue history list
// Return generated dialogue
func GenDialog(history []string) ([]string, string) {
	ans := fmt.Sprintf("generated for %s", history[len(history)-1])
	history = append(history, "doc:"+ans)
	return history, ans
}
