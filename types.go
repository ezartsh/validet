package validet

import "github.com/tidwall/gjson"

type Int = int

type DataObject = map[string]any

type Lookup = func(k string) gjson.Result

type SchemaObject = map[string]any

type SchemaSliceObject = []SchemaObject

type RequiredIf struct {
	FieldPath string
	Value     any
}

type RequiredUnless struct {
	FieldPath string
	Value     any
}
