package validet

import (
	"fmt"
	"reflect"

	"github.com/tidwall/gjson"
)

type BooleanErrorMessage struct {
	Required       string
	RequiredIf     string
	RequiredUnless string
	Custom         string
}

type Boolean struct {
	Required       bool
	RequiredIf     *RequiredIf
	RequiredUnless *RequiredUnless
	Custom         func(v bool, look Lookup) error
	Message        BooleanErrorMessage
}

func (s Boolean) isMyTypeOf(schema any) bool {
	return reflect.TypeOf(schema).Kind() == reflect.Struct && reflect.TypeOf(schema) == reflect.TypeOf(Boolean{})
}

func (s Boolean) process(params RuleParams) ([]string, error) {
	// errorBags := params.ErrorBags
	schemaData := params.DataKey.(DataObject)
	return params.Schema.validate(params.OriginalData, params.Key, schemaData[params.Key], params.Option)
	// pathKey := params.PathKey + params.Key
	// if err != nil {
	// 	errorBags.append(pathKey, bags)
	// 	if params.Option.AbortEarly {
	// 		return errors.New("test")
	// 	}
	// }
	// return nil
}

func (s Boolean) validate(source []byte, key string, value any, option Options) ([]string, error) {
	var bags []string
	err := s.assertRequired(key, value, &bags)

	if err != nil {
		return bags, err
	}

	if err = s.assertRequiredIf(source, key, value, &bags); err != nil {
		return bags, err
	}

	if err = s.assertRequiredUnless(source, key, value, &bags); err != nil {
		return bags, err
	}

	if value != nil {

		stringValue, err := s.assertType(key, value, &bags)

		if err != nil {
			return bags, err
		}

		if s.Custom != nil {
			if err := s.assertCustomValidation(s.Custom, source, stringValue, &bags); option.AbortEarly && err != nil {
				return bags, err
			}
		}

	}

	if len(bags) > 0 {
		return bags, StringValidationError
	}

	return bags, nil
}

func (s Boolean) assertType(key string, value any, bags *[]string) (bool, error) {
	var booleanValue bool
	if isBooelanValue(value) {
		booleanValue = value.(bool)
	} else {
		appendErrorBags(
			bags,
			fmt.Sprintf("%s must be type of boolean", key),
			"",
		)
		return false, BooleanValidationError
	}
	return booleanValue, nil
}

func (s Boolean) assertRequired(key string, value any, bags *[]string) error {
	if s.Required {
		if value == nil {
			appendErrorBags(
				bags,
				fmt.Sprintf("%s is required", key),
				s.Message.Required,
			)
			return BooleanValidationError
		}
	}
	return nil
}

func (s Boolean) assertRequiredIf(jsonSource []byte, key string, value any, bags *[]string) error {
	if s.RequiredIf != nil && value == nil {
		comparedValue := gjson.GetBytes(jsonSource, s.RequiredIf.FieldPath)
		if comparedValue.String() == s.RequiredIf.Value {
			appendErrorBags(
				bags,
				fmt.Sprintf("%s is required", key),
				s.Message.RequiredIf,
			)
			return BooleanValidationError
		}
	}
	return nil
}

func (s Boolean) assertRequiredUnless(jsonSource []byte, key string, value any, bags *[]string) error {
	if s.RequiredUnless != nil && value == nil {
		comparedValue := gjson.GetBytes(jsonSource, s.RequiredUnless.FieldPath)
		if comparedValue.String() != s.RequiredUnless.Value {
			appendErrorBags(
				bags,
				fmt.Sprintf("%s is required", key),
				s.Message.RequiredUnless,
			)
			return StringValidationError
		}
	}
	return nil
}

func (s Boolean) assertCustomValidation(fc func(v bool, look Lookup) error, jsonSource []byte, value any, bags *[]string) error {
	err := fc(value.(bool), func(k string) gjson.Result {
		return gjson.GetBytes(jsonSource, k)
	})
	if err != nil {
		appendErrorBags(
			bags,
			err.Error(),
			s.Message.Custom,
		)
		return BooleanValidationError
	}
	return nil
}
