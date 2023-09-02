package validet

import (
	"errors"
	"fmt"
)

type SliceErrorMessage struct {
	Required string
}

type SliceValueType interface {
	int | int32 | int64 | float32 | float64 | string
}

type Slice[T SliceValueType] struct {
	Required  bool
	ValueType T
	Message   SliceErrorMessage
}

func (s *Slice[T]) validate(jsonSource string, key string, value any, option Options) ([]string, error) {
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

func (s *Slice[T]) assertType(key string, values []any, bags *[]string) error {
	failed := false

	for _, value := range values {
		if _, ok := value.(T); !ok {
			failed = true
			break
		}
	}
	if failed {
		appendErrorBags(
			bags,
			fmt.Sprintf("%s must be type of %T", key, *new(T)),
			"",
		)
		return StringValidationError
	}
	return nil
}
