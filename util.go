package rigo

import "strings"

func SplitCommand(input string) []string {
	return strings.Split(input, " ")
}
