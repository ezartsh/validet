package validet

type Int = int

type DataObject = map[string]any

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
