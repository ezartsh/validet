package validet

import (
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
)

type SliceObjectErrorMessage struct {
	Required       string
	RequiredIf     string
	RequiredUnless string
	Min            string
	Max            string
}

type SliceObject struct {
	Required       bool
	RequiredIf     *RequiredIf
	RequiredUnless *RequiredUnless
	Min            int
	Max            int
	Item           DataObject
	Message        SliceObjectErrorMessage
}

var SliceObjectValidationError = errors.New("slice object validation failed")

func (s *SliceObject) validate(jsonSource string, key string, value any, option Options) ([]string, error) {
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
		values := value.([]DataObject)

		if len(values) > 0 {

			if err := s.assertMin(key, values, &bags); option.AbortEarly && err != nil {
				return bags, err
			}

			if err := s.assertMax(key, values, &bags); option.AbortEarly && err != nil {
				return bags, err
			}
		}

	}

	if len(bags) > 0 {
		return bags, SliceObjectValidationError
	}

	return bags, nil

}

func (s *SliceObject) assertRequired(key string, value any, bags *[]string) error {
	if s.Required {
		if value == nil {
			appendErrorBags(
				bags,
				fmt.Sprintf("%s is required", key),
				s.Message.Required,
			)
			return SliceObjectValidationError
		}
		values, err := s.assertType(key, value, bags)
		if err != nil {
			return err
		}
		if len(values) == 0 {
			appendErrorBags(
				bags,
				fmt.Sprintf("%s is required", key),
				s.Message.Required,
			)
			return SliceObjectValidationError
		}
	}
	return nil
}

func (s *SliceObject) assertType(key string, value any, bags *[]string) ([]DataObject, error) {
	if values, ok := value.([]DataObject); ok {
		return values, nil
	}
	appendErrorBags(
		bags,
		fmt.Sprintf("%s must be type of data object", key),
		"",
	)
	return []DataObject{}, SliceObjectValidationError
}

func (s *SliceObject) assertRequiredIf(jsonSource string, key string, value any, bags *[]string) error {
	values := value.([]DataObject)
	if s.RequiredIf != nil && (value == nil || len(values) == 0) {
		comparedValue := gjson.Get(jsonSource, s.RequiredIf.FieldPath)
		if comparedValue.String() == s.RequiredIf.Value {
			appendErrorBags(
				bags,
				fmt.Sprintf("%s is required", key),
				s.Message.RequiredIf,
			)
			return SliceObjectValidationError
		}
	}
	return nil
}

func (s *SliceObject) assertRequiredUnless(jsonSource string, key string, value any, bags *[]string) error {
	values := value.([]DataObject)
	if s.RequiredUnless != nil && (value == nil || len(values) == 0) {
		comparedValue := gjson.Get(jsonSource, s.RequiredUnless.FieldPath)
		if comparedValue.String() != s.RequiredUnless.Value {
			appendErrorBags(
				bags,
				fmt.Sprintf("%s is required", key),
				s.Message.RequiredUnless,
			)
			return SliceObjectValidationError
		}
	}
	return nil
}

func (s *SliceObject) assertMin(key string, values []DataObject, bags *[]string) error {
	if s.Min > 0 && len(values) < s.Min {
		appendErrorBags(
			bags,
			fmt.Sprintf("%s must be minimum of %d", key, s.Min),
			s.Message.Min,
		)
		return SliceObjectValidationError
	}
	return nil
}

func (s *SliceObject) assertMax(key string, values []DataObject, bags *[]string) error {
	if s.Max > 0 && len(values) > s.Max {
		appendErrorBags(
			bags,
			fmt.Sprintf("%s must be maximum of %d", key, s.Max),
			s.Message.Max,
		)
		return SliceObjectValidationError
	}
	return nil
}
