package validet

type DataObject = map[string]interface{}

type SchemaObject = map[string]interface{}

type SchemaSliceObject = []SchemaObject

type RequiredIf struct {
	FieldPath string
	Value     any
}
