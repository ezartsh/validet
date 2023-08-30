package validet

import (
	"errors"
	"fmt"
)

type ObjectErrorMessage struct {
	Required string
}

type Object struct {
	Required bool
	Item     map[string]any
	Message  ObjectErrorMessage
}

func (s *Object) validate(jsonSource string, key string, value any, option Options) ([]string, error) {
	var bags []string
	if value == nil {
		bags = append(bags, fmt.Sprintf("%s is required", key))
		return bags, errors.New("validation failed")
	}
	stringValue := value.(map[string]any)
	if s.Required {
		if len(stringValue) == 0 {
			bags = append(bags, fmt.Sprintf("%s is required", key))
			return bags, errors.New("validation failed")
		}
	}

	if len(bags) > 0 {
		return bags, errors.New("validation failed")
	}

	return bags, nil

}
