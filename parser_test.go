package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"runtime/pprof"
	"strings"
	"testing"
	// "github.com/google/go-cmp/cmp"
)

func PrettyPrintMap(m map[string]any) string {
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	return string(b)
}

func ConvertIntsToFloat64(data any) any {
	switch v := data.(type) {
	case map[string]any:
		for k, val := range v {
			v[k] = ConvertIntsToFloat64(val)
		}
	case []any:
		for i, val := range v {
			v[i] = ConvertIntsToFloat64(val)
		}
	case int64: // Ensure conversion
		return float64(v)
	}
	return data
}

func TestParsing(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name: "Simple JSON with 1 level nesting",
			input: `{
				"key1": 123,
				"key2": "value2",
				"key2.5": null,
				"key3": {
					"key4": 456,
					"key5": "value5"
				}
			}`,
		},
		{
			name: "JSON with array",
			input: `{
				"numbers": [1, 2, 3, 4.5],
				"nested": {
					"key": 789
				}
			}`,
		},
		{
			name:  "Empty JSON object",
			input: `{}`,
		},
		{
			name: "JSON with boolean values",
			input: `{
				"trueVal": true,
				"falseVal": false
			}`,
		},
		{
			name: "JSON with mixed types in array",
			input: `{
				"array": [123, "text", null, {"nestedInt": 42}]
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rd := strings.NewReader(tt.input)

			parser := NewParser(rd)
			actual := parser.Parse()
			actual = ConvertIntsToFloat64(actual)

			expected := make(map[string]any)
			json.Unmarshal([]byte(tt.input), &expected)

			if !reflect.DeepEqual(actual, expected) {
				t.Errorf("Test %q failed.\nExpected: %v\nActual:   %v", tt.name, PrettyPrintMap(expected), PrettyPrintMap(actual.(map[string]any)))
			}
			// if diff := cmp.Diff(expected, actual); diff != "" {
			// 	t.Errorf("Mismatch (-expected +actual):\n%s", diff)
			// }
		})
	}
}

var EXAMPLEJSON = `{
	"name": "John",
	"age": 27,
	"cars": [
		{
			"model_name": "Honda 2002",
			"vin": null,
			"years": 23,
			"needs_maintnence": true
		},
		{
			"model_name": "Water 7832",
			"vin": 107810204019401,
			"years": 5.80975e3,
			"needs_maintnence": false
		}
	]
})
`

var FROMTEST = `
{
    "st": 1,
    "sid": 486,
    "tt": "active",
    "gr": 0,
    "uuid": "de305d54-75b4-431b-adb2-eb6b9e546014",
    "ip": "127.0.0.1",
    "ua": "user_agent",
    "tz": -6,
    "v": 1
}
`

func BenchmarkParsing(b *testing.B) {
	b.ReportAllocs()
	b.SetBytes(int64(len(EXAMPLEJSON)))
	b.ResetTimer()
	for range b.N {
		rd := strings.NewReader(EXAMPLEJSON)
		parser := NewParser(rd)
		parser.Parse()
	}
}

func BenchmarkSmallLoad(b *testing.B) {
	b.ReportAllocs()
	b.SetBytes(int64(len(FROMTEST)))
	file, _ := os.Create("mem.pprof")
	pprof.WriteHeapProfile(file)
	b.ResetTimer()
	for range b.N {
		rd := strings.NewReader(FROMTEST)
		parser := NewParser(rd)
		parser.Parse()
	}
}
