package main

import (
	"fmt"
	"log/slog"
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

	lexer := newLexer(inputString)

	slog.SetLogLoggerLevel(slog.LevelDebug)
	for {
		t := lexer.nextToken()
		fmt.Println("out token", t)
		if t == tokenEOF {
			fmt.Println("end of json stream")
			break
		}
	}
}
