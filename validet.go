package validet

import (
	"encoding/json"
	"errors"
	"reflect"
)

type Options struct {
	AbortEarly bool
}

type Validation struct {
	data    DataObject
	schema  SchemaRules
	options Options
}

func (v *Validation) check(b *ErrorBag) {
	errorBags := *b

	jsonData, err := json.Marshal(v.data)

	if err != nil {
		jsonData = []byte{}
	}

	mapSchemas(jsonData, "", "", v.data, v.schema, &errorBags, v.options)
}

func mapSchemas(jsonString []byte, pathKey string, key string, data any, schema any, b *ErrorBag, option Options) {
	errorBags := *b

	if pathKey != "" {
		pathKey = pathKey + "."
	}

	schemaData := data.(DataObject)
	if isSchemaRule(schema) {
		schemaRules := schema.(map[string]Rule)
		for scKey, scRules := range schemaRules {
			mapSchemas(jsonString, pathKey, scKey, schemaData, scRules, &errorBags, option)
			if option.AbortEarly && len(errorBags.Errors) > 0 {
				return
			}
		}
	} else {
		if schemaRule, ok := isRule(schema); ok {
			if schemaRule.isMyTypeOf(schema) {
				schemaRule.process(RuleParams{
					OriginalData: jsonString,
					DataKey:      schemaData,
					PathKey:      pathKey,
					Key:          key,
					Schema:       schemaRule,
					ErrorBags:    &errorBags,
					Option:       option,
				})
			}
		}
	}
}

func validate(d DataObject, schema map[string]Rule, options Options) (ErrorBag, error) {
	var errorBags = NewErrorBags()
	validation := Validation{
		data:    d,
		schema:  schema,
		options: options,
	}
	validation.check(errorBags)

	if len(errorBags.Errors) > 0 {
		return *errorBags, errors.New("error validation inputs.")
	}
	return ErrorBag{}, nil
}

func isSchemaRule(val any) bool {
	if value, ok := val.(map[string]Rule); ok {
		return reflect.TypeOf(value) == reflect.TypeOf(map[string]Rule{})
	}
	return false
}

func isRule(val any) (Rule, bool) {
	var value Rule
	if value, ok := val.(Rule); ok {
		return value, ok
	}
	return value, false
}
