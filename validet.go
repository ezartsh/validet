package validet

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

type Options struct {
	AbortEarly bool
}

type Validation struct {
	data    DataObject
	schema  SchemaObject
	options Options
}

func (v *Validation) validate(b *ErrorBag) {
	errorBags := *b

	jsonString, err := json.Marshal(v.data)

	if err != nil {
		jsonString = []byte{}
	}

	v.mapSchemaObject(string(jsonString), "", "", v.data, v.schema, &errorBags)
}

func (v *Validation) mapSchemaObject(jsonString string, pathKey string, key string, data any, schema any, b *ErrorBag) {
	errorBags := *b

	if pathKey != "" {
		pathKey = pathKey + "."
	}

	if isSchema(schema) {
		schemaObject := schema.(SchemaObject)
		schemaData := data.(DataObject)
		for scKey, scValue := range schemaObject {
			v.mapSchemaObject(jsonString, pathKey, scKey, schemaData, scValue, &errorBags)
			if v.options.AbortEarly && len(errorBags.Errors) > 0 {
				return
			}
		}
	} else {
		if isSchemaOfObject(schema) {
			schemaData := data.(DataObject)
			if scObject, ok := schema.(Object); ok {
				bags, err := scObject.validate(jsonString, key, schemaData[key], v.options)
				if err != nil {
					errorBags.append(pathKey+key, bags)
				} else {
					schemaDataValue := schemaData[key].(DataObject)
					for scObjItemKey, scObjItemValue := range scObject.Item {
						v.mapSchemaObject(jsonString, pathKey+key, scObjItemKey, schemaDataValue, scObjItemValue, &errorBags)
						if v.options.AbortEarly && len(errorBags.Errors) > 0 {
							return
						}
					}
				}
			}
		} else if isSchemaOfSlice(schema) {
			schemaData := data.(DataObject)
			var err error
			var bags []string

			switch reflect.TypeOf(schema) {
			case reflect.TypeOf(Slice[string]{}):
				if scMap, ok := schema.(Slice[string]); ok {
					bags, err = scMap.validate(jsonString, key, schemaData[key], v.options)
				}
			case reflect.TypeOf(Slice[int]{}):
				if scMap, ok := schema.(Slice[int]); ok {
					bags, err = scMap.validate(jsonString, key, schemaData[key], v.options)
				}
			case reflect.TypeOf(Slice[int32]{}):
				if scMap, ok := schema.(Slice[int32]); ok {
					bags, err = scMap.validate(jsonString, key, schemaData[key], v.options)
				}
			case reflect.TypeOf(Slice[int64]{}):
				if scMap, ok := schema.(Slice[int64]); ok {
					bags, err = scMap.validate(jsonString, key, schemaData[key], v.options)
				}
			case reflect.TypeOf(Slice[float32]{}):
				if scMap, ok := schema.(Slice[float32]); ok {
					bags, err = scMap.validate(jsonString, key, schemaData[key], v.options)
				}
			case reflect.TypeOf(Slice[float64]{}):
				if scMap, ok := schema.(Slice[float64]); ok {
					bags, err = scMap.validate(jsonString, key, schemaData[key], v.options)
				}
			}
			if err != nil {
				errorBags.append(key, bags)
				if v.options.AbortEarly {
					return
				}
			}
		} else if isSchemaOfSliceObject(schema) {
			if scSliceObject, ok := schema.(SliceObject); ok {
				schemaData := data.(DataObject)
				bags, err := scSliceObject.validate(jsonString, key, schemaData[key], v.options)
				if err != nil {
					errorBags.append(key, bags)
					if v.options.AbortEarly {
						return
					}
				} else {
					schemaDataValues := schemaData[key].([]DataObject)
					for i, value := range schemaDataValues {
						for scObjItemKey, scObjItemValue := range scSliceObject.Item {
							v.mapSchemaObject(jsonString, pathKey+key+"."+strconv.Itoa(i), scObjItemKey, value, scObjItemValue, &errorBags)
							if v.options.AbortEarly && len(errorBags.Errors) > 0 {
								return
							}
						}
					}
				}
			}
		} else if isSchemaOfString(schema) {
			if scMap, ok := schema.(String); ok {
				schemaData := data.(DataObject)
				bags, err := scMap.validate(jsonString, key, schemaData[key], v.options)
				pathKey = pathKey + key
				if err != nil {
					errorBags.append(pathKey, bags)
					if v.options.AbortEarly {
						return
					}
				}
			}
		} else if isSchemaOfNumeric(schema) {
			schemaData := data.(DataObject)
			var err error
			var bags []string

			switch reflect.TypeOf(schema) {
			case reflect.TypeOf(Numeric[int]{}):
				if scMap, ok := schema.(Numeric[int]); ok {
					bags, err = scMap.validate(jsonString, key, schemaData[key], v.options)
				}
			case reflect.TypeOf(Numeric[int32]{}):
				if scMap, ok := schema.(Numeric[int32]); ok {
					bags, err = scMap.validate(jsonString, key, schemaData[key], v.options)
				}
			case reflect.TypeOf(Numeric[int64]{}):
				if scMap, ok := schema.(Numeric[int64]); ok {
					bags, err = scMap.validate(jsonString, key, schemaData[key], v.options)
				}
			case reflect.TypeOf(Numeric[float32]{}):
				if scMap, ok := schema.(Numeric[float32]); ok {
					bags, err = scMap.validate(jsonString, key, schemaData[key], v.options)
				}
			case reflect.TypeOf(Numeric[float64]{}):
				if scMap, ok := schema.(Numeric[float64]); ok {
					bags, err = scMap.validate(jsonString, key, schemaData[key], v.options)
				}
			}
			pathKey = pathKey + key
			if err != nil {
				errorBags.append(pathKey, bags)
				if v.options.AbortEarly {
					return
				}
			}
		}
	}
}

func Validate(d DataObject, schema SchemaObject, options Options) {
	var errorBags = NewErrorBags()
	validation := Validation{
		data:    d,
		schema:  schema,
		options: options,
	}
	validation.validate(errorBags)

	for key, strings := range errorBags.Errors {
		fmt.Println(key, strings)
	}
}

func isSchemaOfString(val any) bool {
	return reflect.TypeOf(val).Kind() == reflect.Struct && reflect.TypeOf(val) == reflect.TypeOf(String{})
}

func isSchemaOfNumeric(val any) bool {
	return reflect.TypeOf(val).Kind() == reflect.Struct && (reflect.TypeOf(val) == reflect.TypeOf(Numeric[int]{}) ||
		reflect.TypeOf(val) == reflect.TypeOf(Numeric[int32]{}) ||
		reflect.TypeOf(val) == reflect.TypeOf(Numeric[int64]{}) ||
		reflect.TypeOf(val) == reflect.TypeOf(Numeric[float32]{}) ||
		reflect.TypeOf(val) == reflect.TypeOf(Numeric[float64]{}))
}

func isSchemaOfObject(val any) bool {
	return reflect.TypeOf(val).Kind() == reflect.Struct && reflect.TypeOf(val) == reflect.TypeOf(Object{})
}

func isSchemaOfSlice(val any) bool {
	return reflect.TypeOf(val).Kind() == reflect.Struct && (reflect.TypeOf(val) == reflect.TypeOf(Slice[int]{}) ||
		reflect.TypeOf(val) == reflect.TypeOf(Slice[int32]{}) ||
		reflect.TypeOf(val) == reflect.TypeOf(Slice[int64]{}) ||
		reflect.TypeOf(val) == reflect.TypeOf(Slice[float32]{}) ||
		reflect.TypeOf(val) == reflect.TypeOf(Slice[float64]{}) ||
		reflect.TypeOf(val) == reflect.TypeOf(Slice[string]{}))
}

func isSchemaOfSliceObject(val any) bool {
	return reflect.TypeOf(val).Kind() == reflect.Struct && reflect.TypeOf(val) == reflect.TypeOf(SliceObject{})
}

func isSchema(val any) bool {
	if value, ok := val.(SchemaObject); ok {
		return reflect.TypeOf(value) == reflect.TypeOf(SchemaObject{})
	}
	return false
}
