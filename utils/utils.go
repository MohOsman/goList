package utils

import (
	"goList/types"
	"unicode"
)

func IsEmptyList(tasks []types.Task) bool {
	return len(tasks) == 0
}
// mapte include nummers to username 
func IsUsernameValid(username string) bool {
	for _, char := range username {
		if !unicode.IsLetter(char) {
			return false
		}
	}
	return true
}