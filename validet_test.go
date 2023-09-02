package validet

import (
	"testing"
)

func TestValidate(t *testing.T) {
	data := DataObject{
		"name":        "tono",
		"email":       "",
		"description": "",
		"url":         "http://www.ada.com",
		"information": DataObject{
			"age":         1.2432,
			"description": "ada",
			"job": DataObject{
				"level": "",
			},
		},
		"tags": []any{"sts"},
		"items": []DataObject{
			{"title": "ada"},
		},
	}
	schema := SchemaObject{
		"name": String{Required: true, Max: 10, Message: StringErrorMessage{Required: "name dibutuhkan"}},
		"email": String{RequiredUnless: &RequiredUnless{
			FieldPath: "name",
			Value:     "tono",
		}, Email: true},
		"url": String{Required: true, Url: &Url{Https: true}},
		"description": String{
			RequiredIf: &RequiredIf{
				FieldPath: "name",
				Value:     "tono",
			},
		},
		"information": Object{
			Required: true,
			Item: SchemaObject{
				"age": Numeric[float64]{RequiredIf: &RequiredIf{
					FieldPath: "name",
					Value:     "tono",
				}, MinDigits: 5},
				"description": String{Required: true, Max: 10, Regex: "p([a-z]+)ch"},
				"job": Object{
					Required: true,
					Item: SchemaObject{
						"level":       String{Required: true, Max: 10},
						"description": String{Required: true, Max: 10},
					},
				},
			},
		},
		"tags": Slice[Int]{Required: true},
		"items": SliceObject{
			Required: true,
			Item: SchemaObject{
				"title":  String{Required: true, Max: 10},
				"status": String{Required: true, Max: 10},
			},
		},
	}

	Validate(data, schema, Options{})

}
