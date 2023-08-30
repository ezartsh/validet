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
	Item     map[string]interface{}
	Message  ObjectErrorMessage
}

func (s *Object) validate(jsonSource string, key string, value interface{}) ([]string, error) {
	var bags []string
	if value == nil {
		bags = append(bags, fmt.Sprintf("%s is required", key))
		return bags, errors.New("validation failed")
	}
	stringValue := value.(map[string]interface{})
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
