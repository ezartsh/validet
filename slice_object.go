package validet

import (
	"errors"
	"fmt"
)

type SliceObjectErrorMessage struct {
	Required string
}

type SliceObject struct {
	Required bool
	Item     map[string]interface{}
	Message  SliceObjectErrorMessage
}

func (s *SliceObject) validate(jsonSource string, key string, value interface{}) ([]string, error) {
	var bags []string
	if value == nil {
		bags = append(bags, fmt.Sprintf("%s is required", key))
		return bags, errors.New("validation failed")
	}
	values := value.([]map[string]interface{})
	if s.Required {
		if len(values) == 0 {
			bags = append(bags, fmt.Sprintf("%s is required", key))
			return bags, errors.New("validation failed")
		}
	}

	if len(bags) > 0 {
		return bags, errors.New("validation failed")
	}

	return bags, nil

}
