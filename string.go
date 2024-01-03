package validet

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strings"

	"github.com/tidwall/gjson"
)

type StringErrorMessage struct {
	Required       string
	RequiredIf     string
	RequiredUnless string
	Min            string
	Max            string
	Regex          string
	NotRegex       string
	In             string
	NotIn          string
	Email          string
	Alpha          string
	AlphaNumeric   string
	Url            string
}

type String struct {
	Required       bool
	RequiredIf     *RequiredIf
	RequiredUnless *RequiredUnless
	Min            int
	Max            int
	Regex          string
	NotRegex       string
	In             []string
	NotIn          []string
	Email          bool
	Alpha          bool
	AlphaNumeric   bool
	Url            *Url
	Custom         func(v string) error
	Message        StringErrorMessage
}

type Url struct {
	Http  bool
	Https bool
}

const (
	urlHttp  string = "http"
	urlHttps        = "https"
)

func (s String) isMyTypeOf(schema any) bool {
	return reflect.TypeOf(schema).Kind() == reflect.Struct && reflect.TypeOf(schema) == reflect.TypeOf(String{})
}

func (s String) process(params RuleParams) error {
	errorBags := params.ErrorBags
	schemaData := params.DataKey.(DataObject)
	bags, err := params.Schema.validate(params.OriginalData, params.Key, schemaData[params.Key], params.Option)
	pathKey := params.PathKey + params.Key
	if err != nil {
		errorBags.append(pathKey, bags)
		if params.Option.AbortEarly {
			return errors.New("test")
		}
	}
	return nil
}

func (s String) validate(source []byte, key string, value any, option Options) ([]string, error) {
	var bags []string
	err := s.assertRequired(key, value, &bags)

	if err != nil {
		return bags, err
	}

	if err = s.assertRequiredIf(source, key, value, &bags); err != nil {
		return bags, err
	}

	if err = s.assertRequiredUnless(source, key, value, &bags); err != nil {
		return bags, err
	}

	if value != nil {

		stringValue, err := s.assertType(key, value, &bags)

		if err != nil {
			return bags, err
		}

		if err := s.assertMin(key, stringValue, &bags); option.AbortEarly && err != nil {
			return bags, err
		}

		if err := s.assertMax(key, stringValue, &bags); option.AbortEarly && err != nil {
			return bags, err
		}

		if err := s.assertRegex(key, stringValue, &bags); option.AbortEarly && err != nil {
			return bags, err
		}

		if err := s.assertNotRegex(key, stringValue, &bags); option.AbortEarly && err != nil {
			return bags, err
		}

		if err := s.assertIn(key, stringValue, &bags); option.AbortEarly && err != nil {
			return bags, err
		}

		if err := s.assertNotIn(key, stringValue, &bags); option.AbortEarly && err != nil {
			return bags, err
		}

		if err := s.assertEmail(key, stringValue, &bags); option.AbortEarly && err != nil {
			return bags, err
		}

		if err := s.assertAlpha(key, stringValue, &bags); option.AbortEarly && err != nil {
			return bags, err
		}

		if err := s.assertAlphaNumeric(key, stringValue, &bags); option.AbortEarly && err != nil {
			return bags, err
		}

		if err := s.assertUrl(key, stringValue, &bags); option.AbortEarly && err != nil {
			return bags, err
		}

		if s.Custom != nil {
			if err := s.assertCustomValidation(s.Custom, stringValue, &bags); option.AbortEarly && err != nil {
				return bags, err
			}
		}

	}

	if len(bags) > 0 {
		return bags, StringValidationError
	}

	return bags, nil
}

func (s String) assertType(key string, value any, bags *[]string) (string, error) {
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

func (s String) assertRequired(key string, value any, bags *[]string) error {
	if s.Required {
		if value == nil {
			appendErrorBags(
				bags,
				fmt.Sprintf("%s is required", key),
				s.Message.Required,
			)
			return StringValidationError
		}
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

func (s String) assertRequiredIf(jsonSource []byte, key string, value any, bags *[]string) error {
	if s.RequiredIf != nil && (value == nil || (isStringValue(value) && stringLength(value) == 0)) {
		comparedValue := gjson.GetBytes(jsonSource, s.RequiredIf.FieldPath)
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

func (s String) assertRequiredUnless(jsonSource []byte, key string, value any, bags *[]string) error {
	if s.RequiredUnless != nil && (value == nil || (isStringValue(value) && stringLength(value) == 0)) {
		comparedValue := gjson.GetBytes(jsonSource, s.RequiredUnless.FieldPath)
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

func (s String) assertMin(key string, value string, bags *[]string) error {
	if s.Min > 0 && stringLength(value) < s.Min {
		appendErrorBags(
			bags,
			fmt.Sprintf("%s must be minimum of %d character(s)", key, s.Min),
			s.Message.Min,
		)
		return StringValidationError
	}
	return nil
}

func (s String) assertMax(key string, value string, bags *[]string) error {
	if s.Max > 0 && stringLength(value) > s.Max {
		appendErrorBags(
			bags,
			fmt.Sprintf("%s must be maximum of %d character(s)", key, s.Max),
			s.Message.Max,
		)
		return StringValidationError
	}
	return nil
}

func (s String) assertRegex(key string, value string, bags *[]string) error {
	regx, err := regexp.Compile(s.Regex)
	if s.Regex != "" && stringLength(value) > 0 && (err != nil || !regx.MatchString(value)) {
		appendErrorBags(
			bags,
			fmt.Sprintf("%s is not a valid format", key),
			s.Message.Regex,
		)
		return StringValidationError
	}
	return nil
}

func (s String) assertNotRegex(key string, value string, bags *[]string) error {
	regx, err := regexp.Compile(s.Regex)
	if s.NotRegex != "" && stringLength(value) > 0 && (err != nil || regx.MatchString(value)) {
		appendErrorBags(
			bags,
			fmt.Sprintf("%s is not a valid format", key),
			s.Message.NotRegex,
		)
		return StringValidationError
	}
	return nil
}

func (s String) assertIn(key string, value string, bags *[]string) error {
	if len(s.In) > 0 && stringLength(value) > 0 && !slices.Contains(s.In, value) {
		appendErrorBags(
			bags,
			fmt.Sprintf("%s must in %s", key, strings.Join(s.In, ", ")),
			s.Message.In,
		)
		return StringValidationError
	}
	return nil
}

func (s String) assertNotIn(key string, value string, bags *[]string) error {
	if len(s.NotIn) > 0 && stringLength(value) > 0 && slices.Contains(s.NotIn, value) {
		appendErrorBags(
			bags,
			fmt.Sprintf("%s must not in %s", key, strings.Join(s.NotIn, ", ")),
			s.Message.NotIn,
		)
		return StringValidationError
	}
	return nil
}

func (s String) assertEmail(key string, value string, bags *[]string) error {
	if s.Email && stringLength(value) > 0 {
		regx, err := regexp.Compile(`^([a-zA-Z0-9._%-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,})$`)
		if err != nil || !regx.MatchString(value) {
			appendErrorBags(
				bags,
				fmt.Sprintf("%s is not a valid email", key),
				s.Message.NotRegex,
			)
			return StringValidationError
		}
	}
	return nil
}

func (s String) assertAlpha(key string, value string, bags *[]string) error {
	if s.Alpha && stringLength(value) > 0 {
		regx, err := regexp.Compile(`^[a-zA-Z]+$`)
		if err != nil || !regx.MatchString(value) {
			appendErrorBags(
				bags,
				fmt.Sprintf("%s is not an alphabetic value", key),
				s.Message.NotRegex,
			)
			return StringValidationError
		}
	}
	return nil
}

func (s String) assertAlphaNumeric(key string, value string, bags *[]string) error {
	if s.AlphaNumeric && stringLength(value) > 0 {
		regx, err := regexp.Compile(`^[a-zA-Z0-9]+$`)
		if err != nil || !regx.MatchString(value) {
			appendErrorBags(
				bags,
				fmt.Sprintf("%s is not an alphabetic number value", key),
				s.Message.NotRegex,
			)
			return StringValidationError
		}
	}
	return nil
}

func (s String) assertUrl(key string, value string, bags *[]string) error {
	if s.Url != nil && stringLength(value) > 0 {
		var prefix []string
		if s.Url.Http {
			prefix = append(prefix, urlHttp)
		}
		if s.Url.Https {
			prefix = append(prefix, urlHttps)
		}
		if len(prefix) == 0 {
			prefix = []string{"http", "https"}
		}
		expression := `^((` + strings.Join(prefix, "|") + `):\/\/)[-a-zA-Z0-9@:%._\\+~#?&\/=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%._\\+~#?&\/=]*)$`
		regx, err := regexp.Compile(expression)
		if err != nil || !regx.MatchString(value) {
			appendErrorBags(
				bags,
				fmt.Sprintf("%s is not a valid url", key),
				s.Message.Url,
			)
			return StringValidationError
		}
	}
	return nil
}

func (s String) assertCustomValidation(fc func(v string) error, value any, bags *[]string) error {
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
