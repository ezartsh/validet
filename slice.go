package validet

import (
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
)

type SliceErrorMessage struct {
	Required       string
	RequiredIf     string
	RequiredUnless string
	Min            string
	Max            string
}

type SliceValueType interface {
	int | int32 | int64 | float32 | float64 | string
}

type Slice[T SliceValueType] struct {
	Required       bool
	RequiredIf     *RequiredIf
	RequiredUnless *RequiredUnless
	ValueType      T
	Min            int
	Max            int
	Message        SliceErrorMessage
}

func (s *Slice[T]) validate(jsonSource string, key string, value any, option Options) ([]string, error) {
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
		}

	}

	if len(bags) > 0 {
		return bags, errors.New("validation failed")
	}

	return bags, nil

}

func (s *Slice[T]) assertRequired(key string, value any, bags *[]string) error {
	if s.Required {
		if value == nil {
			appendErrorBags(
				bags,
				fmt.Sprintf("%s is required", key),
				s.Message.Required,
			)
			return StringValidationError
		}
		values := value.([]any)
		if len(values) == 0 {
			appendErrorBags(
				bags,
				fmt.Sprintf("%s is required", key),
				s.Message.Required,
			)
			return StringValidationError
		}
	}
	return nil
}

func (s *Slice[T]) assertType(key string, values []any, bags *[]string) ([]T, error) {
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
		return []T{}, StringValidationError
	}

	return parsedValues, nil
}

func (s *Slice[T]) assertRequiredIf(jsonSource string, key string, value any, bags *[]string) error {
	values := value.([]any)
	if s.RequiredIf != nil && (value == nil || len(values) == 0) {
		comparedValue := gjson.Get(jsonSource, s.RequiredIf.FieldPath)
		if comparedValue.String() == s.RequiredIf.Value {
			appendErrorBags(
				bags,
				fmt.Sprintf("%s is required", key),
				s.Message.RequiredIf,
			)
			return StringValidationError
		}
	}
	return nil
}

func (s *Slice[T]) assertRequiredUnless(jsonSource string, key string, value any, bags *[]string) error {
	values := value.([]any)
	if s.RequiredUnless != nil && (value == nil || len(values) == 0) {
		comparedValue := gjson.Get(jsonSource, s.RequiredUnless.FieldPath)
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

func (s *Slice[T]) assertMin(key string, values []T, bags *[]string) error {
	if s.Min > 0 && len(values) < s.Min {
		appendErrorBags(
			bags,
			fmt.Sprintf("%s must be minimum of %d", key, s.Min),
			s.Message.Min,
		)
		return StringValidationError
	}
	return nil
}

func (s *Slice[T]) assertMax(key string, values []T, bags *[]string) error {
	if s.Max > 0 && len(values) > s.Max {
		appendErrorBags(
			bags,
			fmt.Sprintf("%s must be maximum of %d", key, s.Max),
			s.Message.Max,
		)
		return StringValidationError
	}
	return nil
}
