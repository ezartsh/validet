package validet

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/tidwall/gjson"
)

type ObjectErrorMessage struct {
	Required       string
	RequiredIf     string
	RequiredUnless string
}

type Object struct {
	Required       bool
	RequiredIf     *RequiredIf
	RequiredUnless *RequiredUnless
	Item           DataObject
	Message        ObjectErrorMessage
}

func (s Object) isMyTypeOf(schema any) bool {
	return reflect.TypeOf(schema).Kind() == reflect.Struct && reflect.TypeOf(schema) == reflect.TypeOf(Object{})
}

func (s Object) process(params RuleParams) error {
	schemaData := params.DataKey.(DataObject)
	var err error
	var bags []string

	errorBags := *params.ErrorBags
	schema := params.Schema
	originalData := params.OriginalData
	key := params.Key
	options := params.Option

	if scObject, ok := schema.(Object); ok {
		bags, err := scObject.validate(originalData, key, schemaData[key], options)
		if err != nil {
			errorBags.append(params.PathKey+key, bags)
		} else {
			schemaDataValue := schemaData[key].(DataObject)
			for scObjItemKey, scObjItemValue := range scObject.Item {
				mapSchemas(originalData, params.PathKey+key, scObjItemKey, schemaDataValue, scObjItemValue, &errorBags, options)
				if options.AbortEarly && len(errorBags.Errors) > 0 {
					return errors.New("error")
				}
			}
		}
	}

	pathKey := params.PathKey + key
	if err != nil {
		params.ErrorBags.append(pathKey, bags)
		if options.AbortEarly {
			return errors.New("error")
		}
	}
	return nil
}

func (s Object) validate(jsonSource []byte, key string, value any, option Options) ([]string, error) {
	var bags []string

	err := s.assertRequired(key, value, &bags)

	if err != nil {
		return bags, err
	}

	if err = s.assertRequiredIf(jsonSource, key, value, &bags); err != nil {
		return bags, err
	}

	if err = s.assertRequiredUnless(jsonSource, key, value, &bags); err != nil {
		return bags, err
	}

	if len(bags) > 0 {
		return bags, ObjectValidationError
	}

	return bags, nil

}

func (s Object) assertRequired(key string, value any, bags *[]string) error {
	if s.Required {
		if value == nil {
			appendErrorBags(
				bags,
				fmt.Sprintf("%s is required", key),
				s.Message.Required,
			)
			return ObjectValidationError
		}
		if values, ok := value.(DataObject); ok {
			if len(values) == 0 {
				appendErrorBags(
					bags,
					fmt.Sprintf("%s is required", key),
					s.Message.Required,
				)
				return ObjectValidationError
			}
		} else {
			appendErrorBags(
				bags,
				fmt.Sprintf("%s must be type of %T", key, *new(DataObject)),
				s.Message.Required,
			)
			return ObjectValidationError
		}
	}
	return nil
}

func (s Object) assertRequiredIf(jsonSource []byte, key string, value any, bags *[]string) error {
	values := value.(DataObject)
	if s.RequiredIf != nil && (value == nil || len(values) == 0) {
		comparedValue := gjson.GetBytes(jsonSource, s.RequiredIf.FieldPath)
		if comparedValue.String() == s.RequiredIf.Value {
			appendErrorBags(
				bags,
				fmt.Sprintf("%s is required", key),
				s.Message.RequiredIf,
			)
			return ObjectValidationError
		}
	}
	return nil
}

func (s Object) assertRequiredUnless(jsonSource []byte, key string, value any, bags *[]string) error {
	values := value.(DataObject)
	if s.RequiredUnless != nil && (value == nil || len(values) == 0) {
		comparedValue := gjson.GetBytes(jsonSource, s.RequiredUnless.FieldPath)
		if comparedValue.String() != s.RequiredUnless.Value {
			appendErrorBags(
				bags,
				fmt.Sprintf("%s is required", key),
				s.Message.RequiredUnless,
			)
			return ObjectValidationError
		}
	}
	return nil
}
