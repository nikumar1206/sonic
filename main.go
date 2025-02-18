package main

import (
	"fmt"
	"log/slog"
	"strings"
)

func main() {
	// inputString := []byte(`{
	// 	"key1": 123,
	// 	"key2": "value2",
	// 	"key2.5": null,
	// 	"key3": {
	// 		"key4": 456,
	// 		"key5": "value5"
	// 		}
	// 	}`)

	inputString2 := []byte(`{
		"key1" 1e23
		"key2" "value2",
		"key2.5" null,
		"key3" ["foo", "bar", 3, null]
		}`)

	rd := strings.NewReader(string(inputString2))

	parser := NewParser(rd)

	data := parser.Parse()

	fmt.Println("data here", data)
	slog.SetLogLoggerLevel(slog.LevelDebug)
}
