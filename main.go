package main

import (
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

	parser := NewParser(rd)

	data := make(map[string]any)
	parser.Parse(data)
	slog.SetLogLoggerLevel(slog.LevelDebug)
}
