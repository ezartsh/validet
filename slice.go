package validet

import (
	"fmt"
	"reflect"

	"github.com/tidwall/gjson"
)

type SliceErrorMessage struct {
	Required       string
	RequiredIf     string
	RequiredUnless string
	Min            string
	Max            string
	Custom         string
}

type SliceValueType interface {
	int | int32 | int64 | uint | uint32 | uint64 | float32 | float64 | string
}

type Slice[T SliceValueType] struct {
	Required       bool
	RequiredIf     *RequiredIf
	RequiredUnless *RequiredUnless
	Min            int
	Max            int
	Custom         func(v []T, look Lookup) error
	Message        SliceErrorMessage
}

func (s Slice[T]) isMyTypeOf(schema any) bool {
	return reflect.TypeOf(schema).Kind() == reflect.Struct && (reflect.TypeOf(schema) == reflect.TypeOf(Slice[int]{}) ||
		reflect.TypeOf(schema) == reflect.TypeOf(Slice[int32]{}) ||
		reflect.TypeOf(schema) == reflect.TypeOf(Slice[int64]{}) ||
		reflect.TypeOf(schema) == reflect.TypeOf(Slice[uint]{}) ||
		reflect.TypeOf(schema) == reflect.TypeOf(Slice[uint32]{}) ||
		reflect.TypeOf(schema) == reflect.TypeOf(Slice[uint64]{}) ||
		reflect.TypeOf(schema) == reflect.TypeOf(Slice[float32]{}) ||
		reflect.TypeOf(schema) == reflect.TypeOf(Slice[float64]{}) ||
		reflect.TypeOf(schema) == reflect.TypeOf(Slice[string]{}))
}

func (s Slice[T]) process(params RuleParams) ([]string, error) {
	schemaData := params.DataKey.(DataObject)
	var err error
	var bags []string

	schema := params.Schema
	originalData := params.OriginalData
	key := params.Key
	options := params.Option

	switch reflect.TypeOf(schema) {
	case reflect.TypeOf(Slice[string]{}):
		if scMap, ok := schema.(Slice[string]); ok {
			bags, err = scMap.validate(originalData, key, schemaData[key], options)
		}
	case reflect.TypeOf(Slice[int]{}):
		if scMap, ok := schema.(Slice[int]); ok {
			bags, err = scMap.validate(originalData, key, schemaData[key], options)
		}
	case reflect.TypeOf(Slice[int32]{}):
		if scMap, ok := schema.(Slice[int32]); ok {
			bags, err = scMap.validate(originalData, key, schemaData[key], options)
		}
	case reflect.TypeOf(Slice[int64]{}):
		if scMap, ok := schema.(Slice[int64]); ok {
			bags, err = scMap.validate(originalData, key, schemaData[key], options)
		}
	case reflect.TypeOf(Slice[uint]{}):
		if scMap, ok := schema.(Slice[uint]); ok {
			bags, err = scMap.validate(originalData, key, schemaData[key], options)
		}
	case reflect.TypeOf(Slice[uint32]{}):
		if scMap, ok := schema.(Slice[uint32]); ok {
			bags, err = scMap.validate(originalData, key, schemaData[key], options)
		}
	case reflect.TypeOf(Slice[uint64]{}):
		if scMap, ok := schema.(Slice[uint64]); ok {
			bags, err = scMap.validate(originalData, key, schemaData[key], options)
		}
	case reflect.TypeOf(Slice[float32]{}):
		if scMap, ok := schema.(Slice[float32]); ok {
			bags, err = scMap.validate(originalData, key, schemaData[key], options)
		}
	case reflect.TypeOf(Slice[float64]{}):
		if scMap, ok := schema.(Slice[float64]); ok {
			bags, err = scMap.validate(originalData, key, schemaData[key], options)
		}
	}

	return bags, err

	// pathKey := params.PathKey + key
	// if err != nil {
	// 	params.ErrorBags.append(pathKey, bags)
	// 	if options.AbortEarly {
	// 		return errors.New("error")
	// 	}
	// }
	// return []string{}, nil
}

func (s Slice[T]) validate(jsonSource []byte, key string, value any, option Options) ([]string, error) {
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

	if value != nil {
		values := value.([]any)

		if len(values) > 0 {
			parsedValue, err := s.assertType(key, values, &bags)

			if err != nil {
				return bags, err
			}

			if err := s.assertMin(key, parsedValue, &bags); option.AbortEarly && err != nil {
				return bags, err
			}

			if err := s.assertMax(key, parsedValue, &bags); option.AbortEarly && err != nil {
				return bags, err
			}

			if s.Custom != nil {
				if err := s.assertCustomValidation(s.Custom, jsonSource, parsedValue, &bags); option.AbortEarly && err != nil {
					return bags, err
				}
			}
		}

	}

	if len(bags) > 0 {
		return bags, SliceValidationError
	}

	return bags, nil

}

func (s Slice[T]) assertType(key string, values []any, bags *[]string) ([]T, error) {
	failed := false
	var parsedValues []T
	for _, value := range values {
		if parseValue, ok := value.(T); ok {
			parsedValues = append(parsedValues, parseValue)
		} else {
			failed = true
		}
	}
	if failed {
		appendErrorBags(
			bags,
			fmt.Sprintf("%s must be slice of type %T", key, *new(T)),
			"",
		)
		return []T{}, SliceValidationError
	}

	return parsedValues, nil
}

func (s Slice[T]) assertRequired(key string, value any, bags *[]string) error {
	if s.Required {
		if value == nil {
			appendErrorBags(
				bags,
				fmt.Sprintf("%s is required", key),
				s.Message.Required,
			)
			return SliceValidationError
		}
		if values, ok := value.([]any); ok {
			if len(values) == 0 {
				appendErrorBags(
					bags,
					fmt.Sprintf("%s is required", key),
					s.Message.Required,
				)
				return SliceValidationError
			}
		} else {
			appendErrorBags(
				bags,
				fmt.Sprintf("%s must be slice of type %T", key, *new(T)),
				s.Message.Required,
			)
			return SliceValidationError
		}
	}
	return nil
}

func (s Slice[T]) assertRequiredIf(jsonSource []byte, key string, value any, bags *[]string) error {
	values := value.([]any)
	if s.RequiredIf != nil && (value == nil || len(values) == 0) {
		comparedValue := gjson.GetBytes(jsonSource, s.RequiredIf.FieldPath)
		if comparedValue.String() == s.RequiredIf.Value {
			appendErrorBags(
				bags,
				fmt.Sprintf("%s is required", key),
				s.Message.RequiredIf,
			)
			return SliceValidationError
		}
	}
	return nil
}

func (s Slice[T]) assertRequiredUnless(jsonSource []byte, key string, value any, bags *[]string) error {
	values := value.([]any)
	if s.RequiredUnless != nil && (value == nil || len(values) == 0) {
		comparedValue := gjson.GetBytes(jsonSource, s.RequiredUnless.FieldPath)
		if comparedValue.String() != s.RequiredUnless.Value {
			appendErrorBags(
				bags,
				fmt.Sprintf("%s is required", key),
				s.Message.RequiredUnless,
			)
			return SliceValidationError
		}
	}
	return nil
}

func (s Slice[T]) assertMin(key string, values []T, bags *[]string) error {
	if s.Min > 0 && len(values) < s.Min {
		appendErrorBags(
			bags,
			fmt.Sprintf("%s must be minimum of %d", key, s.Min),
			s.Message.Min,
		)
		return SliceValidationError
	}
	return nil
}

func (s Slice[T]) assertMax(key string, values []T, bags *[]string) error {
	if s.Max > 0 && len(values) > s.Max {
		appendErrorBags(
			bags,
			fmt.Sprintf("%s must be maximum of %d", key, s.Max),
			s.Message.Max,
		)
		return SliceValidationError
	}
	return nil
}

func (s Slice[T]) assertCustomValidation(fc func(v []T, look Lookup) error, jsonSource []byte, value []T, bags *[]string) error {
	err := fc(value, func(k string) gjson.Result {
		return gjson.GetBytes(jsonSource, k)
	})
	if err != nil {
		appendErrorBags(
			bags,
			err.Error(),
			s.Message.Custom,
		)
		return SliceValidationError
	}
	return nil
}
