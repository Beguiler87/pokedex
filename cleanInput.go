package main

import "strings"

func cleanInput(text string) []string {
	inputs := strings.Fields(text)
	lowered := []string{}
	for i := 0; i < len(inputs); i++ {
		loweredWords := strings.ToLower(inputs[i])
		lowered = append(lowered, loweredWords)
	}
	return lowered
}
