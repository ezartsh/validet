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

func (v *Validation) mapSchemaObject(jsonString string, pathKey string, key string, data interface{}, schema interface{}, b *ErrorBag) {
	errorBags := *b

	if pathKey != "" {
		pathKey = pathKey + "."
	}

	if isSchema(schema) {
		schemaObject := schema.(SchemaObject)
		schemaData := data.(DataObject)
		for scKey, scValue := range schemaObject {
			v.mapSchemaObject(jsonString, pathKey, scKey, schemaData, scValue, &errorBags)
		}
	} else {
		if isSchemaOfObject(schema) {
			schemaData := data.(DataObject)
			if scObject, ok := schema.(Object); ok {
				bags, err := scObject.validate(jsonString, key, schemaData[key])
				if err != nil {
					errorBags.append(pathKey+key, bags)
				} else {
					schemaDataValue := schemaData[key].(DataObject)
					for scObjItemKey, scObjItemValue := range scObject.Item {
						v.mapSchemaObject(jsonString, pathKey+key, scObjItemKey, schemaDataValue, scObjItemValue, &errorBags)
					}
				}
			}
		} else if isSchemaOfSlice(schema) {
			if scMap, ok := schema.(Slice); ok {
				schemaData := data.(DataObject)
				bags, err := scMap.validate(jsonString, key, schemaData[key])
				if err != nil {
					errorBags.append(key, bags)
				}
			}
		} else if isSchemaOfSliceObject(schema) {
			if scSliceObject, ok := schema.(SliceObject); ok {
				schemaData := data.(DataObject)
				bags, err := scSliceObject.validate(jsonString, key, schemaData[key])
				if err != nil {
					errorBags.append(key, bags)
				} else {
					schemaDataValues := schemaData[key].([]DataObject)
					for i, value := range schemaDataValues {
						for scObjItemKey, scObjItemValue := range scSliceObject.Item {
							v.mapSchemaObject(jsonString, pathKey+key+"."+strconv.Itoa(i), scObjItemKey, value, scObjItemValue, &errorBags)
						}
					}
				}
			}
		} else if isSchemaOfString(schema) {
			if scMap, ok := schema.(String); ok {
				schemaData := data.(DataObject)
				bags, err := scMap.validate(jsonString, key, schemaData[key])
				pathKey = pathKey + key
				if err != nil {
					errorBags.append(pathKey, bags)
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

func isSchemaOfString(val interface{}) bool {
	return reflect.TypeOf(val).Kind() == reflect.Struct && reflect.TypeOf(val) == reflect.TypeOf(String{})
}

func isSchemaOfObject(val interface{}) bool {
	return reflect.TypeOf(val).Kind() == reflect.Struct && reflect.TypeOf(val) == reflect.TypeOf(Object{})
}

func isSchemaOfSlice(val interface{}) bool {
	return reflect.TypeOf(val).Kind() == reflect.Struct && reflect.TypeOf(val) == reflect.TypeOf(Slice{})
}

func isSchemaOfSliceObject(val interface{}) bool {
	return reflect.TypeOf(val).Kind() == reflect.Struct && reflect.TypeOf(val) == reflect.TypeOf(SliceObject{})
}

func isSchema(val interface{}) bool {
	if value, ok := val.(SchemaObject); ok {
		return reflect.TypeOf(value) == reflect.TypeOf(SchemaObject{})
	}
	return false
}
