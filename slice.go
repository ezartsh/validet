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

func (s *Slice) validate(jsonSource string, key string, value interface{}) ([]string, error) {
	var bags []string
	if value == nil {
		bags = append(bags, fmt.Sprintf("%s is required", key))
		return bags, errors.New("validation failed")
	}
	values := value.([]interface{})
	if s.Required {
		if len(values) == 0 {
			bags = append(bags, fmt.Sprintf("%s is required", key))
			return bags, errors.New("validation failed")
		}
	}

	s.assertType(key, values, &bags)

	if len(bags) > 0 {
		return bags, errors.New("validation failed")
	}

	return bags, nil

}

func (s *Slice) assertType(key string, values []interface{}, bags *[]string) error {
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
