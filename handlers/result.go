package handlers

import "strings"

var dummyWords = []string{
	"sun", "moon", "cloud", "river", "forest",
	"code", "data", "API", "server", "function",
	"goroutine", "channel", "request", "response", "handler",
	"memory", "disk", "CPU", "logic", "database",
}

func Result(wordsUsed int) string {
	var selectedWords []string
	for i := 0; i < wordsUsed; i++ {
		selectedWords = append(selectedWords, dummyWords[i%len(dummyWords)])
	}
	result := strings.Join(selectedWords, " ")
	return result
}
