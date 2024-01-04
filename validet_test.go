package validet

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

type CustomString struct {
	Required bool
}

func (cs CustomString) validate(source []byte, key string, value any, option Options) ([]string, error) {
	return []string{"jangan kosong"}, errors.New("error custom")
}

func (cs CustomString) isMyTypeOf(schema any) bool {
	return reflect.TypeOf(schema).Kind() == reflect.Struct && reflect.TypeOf(schema) == reflect.TypeOf(CustomString{})
}

func (cs CustomString) process(params RuleParams) ([]string, error) {
	schemaData := params.DataKey.(DataObject)
	return params.Schema.validate(params.OriginalData, params.Key, schemaData[params.Key], params.Option)
}

func TestValidate(t *testing.T) {
	data := DataObject{
		"name":        "tono",
		"new_name":    false,
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
		"tags": []any{1},
		"items": []DataObject{
			{"titles": "ada", "collections": []DataObject{{"title": "ada"}}},
		},
	}
	schema := NewSchema(
		data,
		map[string]Rule{
			"name":        String{Required: true, Min: 10, Message: StringErrorMessage{Required: "name dibutuhkan"}},
			"new_name":    Boolean{Required: true},
			"custom_name": CustomString{Required: true},
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
			"store": Numeric[uint]{Required: true},
			"information": Object{
				Required: true,
				Item: SchemaObject{
					"age": Numeric[float64]{RequiredIf: &RequiredIf{
						FieldPath: "name",
						Value:     "tono",
					}, MinDigits: 5},
					"description": String{Required: true, Max: 2, Regex: "p([a-z]+)ch"},
					"job": Object{
						Required: true,
						Item: SchemaObject{
							"level":       String{Required: true, Max: 10},
							"description": String{Required: true, Max: 10},
						},
					},
				},
			},
			"tags": Slice[Int]{Required: true, Min: 2},
			"items": SliceObject{
				Required: true,
				Item: SchemaObject{
					"titles": Slice[string]{Required: true, Min: 2},
					"collections": SliceObject{
						Required: true,
						Item: SchemaObject{
							"title":       String{Required: true, Max: 10},
							"description": String{Required: true, Max: 10},
						},
					},
				},
			},
		},
		Options{},
	)

	bags, _ := schema.Validate()

	fmt.Println(bags)
}
