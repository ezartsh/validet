package validet

import (
	"errors"
	"fmt"
	"reflect"
)

type SliceErrorMessage struct {
	Required string
}

type Slice struct {
	Required  bool
	ValueType string
	Message   SliceErrorMessage
}

func (s *Slice) validate(jsonSource string, key string, value any, option Options) ([]string, error) {
	var bags []string
	if value == nil {
		bags = append(bags, fmt.Sprintf("%s is required", key))
		return bags, errors.New("validation failed")
	}
	values := value.([]any)
	if s.Required {
		if len(values) == 0 {
			bags = append(bags, fmt.Sprintf("%s is required", key))
			return bags, errors.New("validation failed")
		}
	}

	if err := s.assertType(key, values, &bags); option.AbortEarly && err != nil {
		return bags, err
	}

	if len(bags) > 0 {
		return bags, errors.New("validation failed")
	}

	return bags, nil

}

func (s *Slice) assertType(key string, values []any, bags *[]string) error {
	failed := false

	if s.ValueType != "" {
		for _, value := range values {
			if s.ValueType == "int" {
				if reflect.TypeOf(value) != reflect.TypeOf(1) {
					failed = true
					break
				}
			} else if s.ValueType == "string" {
				if reflect.TypeOf(value) != reflect.TypeOf("") {
					failed = true
					break
				}
			}
		}
	}
	if failed {
		appendErrorBags(
			bags,
			fmt.Sprintf("%s must be type of %s", key, s.ValueType),
			"",
		)
		return StringValidationError
	}
	return nil
}
