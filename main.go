package main

import (
	"fmt"
	"log/slog"
	"strings"
)

func main() {
	inputString := []byte(`{
		"key1": 123,
		"key2": "value2",
		"key2.5": null,
		"key3": {
			"key4": 456,
			"key5": "value5"
			}
		}`)

	rd := strings.NewReader(string(inputString))

	parser := NewParser(rd, "stack")

	data := parser.Parse()

	fmt.Println("data here", data)
	slog.SetLogLoggerLevel(slog.LevelDebug)
}
