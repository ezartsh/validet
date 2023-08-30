package validet

import (
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"reflect"
	"regexp"
	"slices"
	"strings"
)

type StringErrorMessage struct {
	Required   string
	RequiredIf string
	Min        string
	Max        string
	Regex      string
	In         string
	NotIn      string
}

type String struct {
	Required   bool
	RequiredIf *RequiredIf
	Min        int
	Max        int
	Regex      string
	In         []string
	NotIn      []string
	Custom     func(v string) error
	Message    StringErrorMessage
}

var StringValidationError = errors.New("string validation failed")

func (s *String) validate(jsonSource string, key string, value interface{}) ([]string, error) {
	var bags []string
	err := s.assertRequired(key, value, &bags)

	if err != nil {
		return bags, err
	}

	err = s.assertRequiredIf(jsonSource, key, value, &bags)

	if err != nil {
		return bags, err
	}

	if value != nil {

		stringValue, err := s.assertType(key, value, &bags)

		if err != nil {
			return bags, err
		}

		s.assertMin(key, stringValue, &bags)

		s.assertMax(key, stringValue, &bags)

		s.assertRegex(key, stringValue, &bags)

		s.assertIn(key, stringValue, &bags)

		s.assertNotIn(key, stringValue, &bags)

		if s.Custom != nil {
			s.assertCustomValidation(s.Custom, stringValue, &bags)
		}

	}

	if len(bags) > 0 {
		return bags, StringValidationError
	}

	return bags, nil
}

func (s *String) assertType(key string, value interface{}, bags *[]string) (string, error) {
	var stringValue string
	if isStringValue(value) {
		stringValue = value.(string)
	} else {
		appendErrorBags(
			bags,
			fmt.Sprintf("%s must be type of string", key),
			"",
		)
		return "", StringValidationError
	}
	return stringValue, nil
}

func (s *String) assertRequired(key string, value interface{}, bags *[]string) error {
	if value == nil {
		appendErrorBags(
			bags,
			fmt.Sprintf("%s is required", key),
			s.Message.Required,
		)
		return StringValidationError
	}
	if s.Required {
		if isStringValue(value) {
			if stringLength(value) == 0 {
				appendErrorBags(
					bags,
					fmt.Sprintf("%s is required", key),
					s.Message.Required,
				)
				return StringValidationError
			}
		}
	}
	return nil
}

func (s *String) assertRequiredIf(jsonSource string, key string, value interface{}, bags *[]string) error {
	if s.RequiredIf != nil && (value == nil || (isStringValue(value) && stringLength(value) == 0)) {
		value := gjson.Get(jsonSource, s.RequiredIf.FieldPath)
		if value.String() == s.RequiredIf.Value {
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

func (s *String) assertMin(key string, value string, bags *[]string) error {
	if s.Min > 0 && stringLength(value) < s.Min {
		appendErrorBags(
			bags,
			fmt.Sprintf("%s must be minimum of %d", key, s.Min),
			s.Message.Min,
		)
		return StringValidationError
	}
	return nil
}

func (s *String) assertMax(key string, value string, bags *[]string) error {
	if s.Max > 0 && stringLength(value) > s.Max {
		appendErrorBags(
			bags,
			fmt.Sprintf("%s must be maximum of %d", key, s.Max),
			s.Message.Max,
		)
		return StringValidationError
	}
	return nil
}

func (s *String) assertRegex(key string, value string, bags *[]string) error {
	regx, err := regexp.Compile(s.Regex)
	if err != nil || !regx.MatchString(value) {
		appendErrorBags(
			bags,
			fmt.Sprintf("%s value is a not valid format", key),
			s.Message.Regex,
		)
		return StringValidationError
	}
	return nil
}

func (s *String) assertIn(key string, value string, bags *[]string) error {
	if len(s.In) > 0 && !slices.Contains(s.In, value) {
		appendErrorBags(
			bags,
			fmt.Sprintf("%s value must in %s", key, strings.Join(s.In, ", ")),
			s.Message.In,
		)
		return StringValidationError
	}
	return nil
}

func (s *String) assertNotIn(key string, value string, bags *[]string) error {
	if len(s.In) > 0 && slices.Contains(s.In, value) {
		appendErrorBags(
			bags,
			fmt.Sprintf("%s value must not in %s", key, strings.Join(s.In, ", ")),
			s.Message.In,
		)
		return StringValidationError
	}
	return nil
}

func (s *String) assertCustomValidation(fc func(v string) error, value interface{}, bags *[]string) error {
	err := fc(value.(string))
	if err != nil {
		appendErrorBags(
			bags,
			err.Error(),
			s.Message.Max,
		)
		return StringValidationError
	}
	return nil
}

func isStringValue(value interface{}) bool {
	return reflect.TypeOf(value) == reflect.TypeOf("")
}

func stringLength(value interface{}) int {
	if value == nil {
		return 0
	}
	stringValue := value.(string)
	return len([]rune(stringValue))
}

func appendErrorBags(bags *[]string, om string, cm string) {
	errorBags := *bags
	message := om
	if cm != "" {
		message = cm
	}
	*bags = append(errorBags, message)
}
