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

			parser := NewParser(rd, "stack")
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
		parser := NewParser(rd, "stack")
		parser.Parse()
	}
}

func BenchmarkSmallLoad(b *testing.B) {
	b.ReportAllocs()
	b.SetBytes(int64(len(FROMTEST)))
	b.ResetTimer()
	for range b.N {
		rd := strings.NewReader(FROMTEST)
		parser := NewParser(rd, "stack")
		parser.Parse()
	}
}

var mediumFixture = `{
  "person": {
    "id": "d50887ca-a6ce-4e59-b89f-14f0b5d03b03",
    "name": {
      "fullName": "Leonid Bugaev",
      "givenName": "Leonid",
      "familyName": "Bugaev"
    },
    "email": "leonsbox@gmail.com",
    "gender": "male",
    "location": "Saint Petersburg, Saint Petersburg, RU",
    "geo": {
      "city": "Saint Petersburg",
      "state": "Saint Petersburg",
      "country": "Russia",
      "lat": 59.9342802,
      "lng": 30.3350986
    },
    "bio": "Senior engineer at Granify.com",
    "site": "http://flickfaver.com",
    "avatar": "https://d1ts43dypk8bqh.cloudfront.net/v1/avatars/d50887ca-a6ce-4e59-b89f-14f0b5d03b03",
    "employment": {
      "name": "www.latera.ru",
      "title": "Software Engineer",
      "domain": "gmail.com"
    },
    "facebook": {
      "handle": "leonid.bugaev"
    },
    "github": {
      "handle": "buger",
      "id": 14009,
      "avatar": "https://avatars.githubusercontent.com/u/14009?v=3",
      "company": "Granify",
      "blog": "http://leonsbox.com",
      "followers": 95,
      "following": 10
    },
    "twitter": {
      "handle": "flickfaver",
      "id": 77004410,
      "bio": null,
      "followers": 2,
      "following": 1,
      "statuses": 5,
      "favorites": 0,
      "location": "",
      "site": "http://flickfaver.com",
      "avatar": null
    },
    "linkedin": {
      "handle": "in/leonidbugaev"
    },
    "googleplus": {
      "handle": null
    },
    "angellist": {
      "handle": "leonid-bugaev",
      "id": 61541,
      "bio": "Senior engineer at Granify.com",
      "blog": "http://buger.github.com",
      "site": "http://buger.github.com",
      "followers": 41,
      "avatar": "https://d1qb2nb5cznatu.cloudfront.net/users/61541-medium_jpg?1405474390"
    },
    "klout": {
      "handle": null,
      "score": null
    },
    "foursquare": {
      "handle": null
    },
    "aboutme": {
      "handle": "leonid.bugaev",
      "bio": null,
      "avatar": null
    },
    "gravatar": {
      "handle": "buger",
      "urls": [
      ],
      "avatar": "http://1.gravatar.com/avatar/f7c8edd577d13b8930d5522f28123510",
      "avatars": [
        {
          "url": "http://1.gravatar.com/avatar/f7c8edd577d13b8930d5522f28123510",
          "type": "thumbnail"
        }
      ]
    },
    "fuzzy": false
  },
  "company": null
}`

func BenchmarkMediumPayload(b *testing.B) {
	b.ReportAllocs()
	file, _ := os.Create("mem.pprof")
	err := pprof.WriteHeapProfile(file)
	if err != nil {
		panic(err)
	}

	b.SetBytes(int64(len(mediumFixture)))
	b.ResetTimer()
	for range b.N {
		rd := strings.NewReader(mediumFixture)
		parser := NewParser(rd, "stack")
		parser.Parse()
	}
}

func BenchmarkStdLib(b *testing.B) {
	b.ReportAllocs()
	b.SetBytes(int64(len(mediumFixture)))
	b.ResetTimer()
	for range b.N {
		v := make(map[string]any)
		json.Unmarshal([]byte(mediumFixture), &v)
	}
}
