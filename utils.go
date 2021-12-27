package main

import "strings"

func getCasesOfString(input string) (string, string) {
	capitalizedInput := strings.Title(input)
	firstChar := input[0]
	lowerCaseChar := strings.ToLower(string(firstChar))
	uncapitalizedInput := lowerCaseChar + string(input[1:])
	return capitalizedInput, uncapitalizedInput
}
