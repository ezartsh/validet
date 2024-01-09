package validet

import (
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

type NumericErrorMessage struct {
	Required       string
	RequiredIf     string
	RequiredUnless string
	Min            string
	Max            string
	MinDigits      string
	MaxDigits      string
	Regex          string
	NotRegex       string
	In             string
	NotIn          string
	Custom         string
}

type NumericValue interface {
	int | int32 | int64 | uint | uint32 | uint64 | float32 | float64
}

type Numeric[NT NumericValue] struct {
	Required       bool
	RequiredIf     *RequiredIf
	RequiredUnless *RequiredUnless
	Min            int
	Max            int
	MinDigits      int
	MaxDigits      int
	Regex          string
	NotRegex       string
	In             []NT
	NotIn          []NT
	Custom         func(v NT, params RuleParams, look Lookup) error
	Message        NumericErrorMessage
}

func (s Numeric[NT]) isMyTypeOf(schema any) bool {
	return reflect.TypeOf(schema).Kind() == reflect.Struct &&
		(reflect.TypeOf(schema) == reflect.TypeOf(Numeric[int]{}) ||
			reflect.TypeOf(schema) == reflect.TypeOf(Numeric[int32]{}) ||
			reflect.TypeOf(schema) == reflect.TypeOf(Numeric[int64]{}) ||
			reflect.TypeOf(schema) == reflect.TypeOf(Numeric[uint]{}) ||
			reflect.TypeOf(schema) == reflect.TypeOf(Numeric[uint32]{}) ||
			reflect.TypeOf(schema) == reflect.TypeOf(Numeric[uint64]{}) ||
			reflect.TypeOf(schema) == reflect.TypeOf(Numeric[float32]{}) ||
			reflect.TypeOf(schema) == reflect.TypeOf(Numeric[float64]{}))
}

func (s Numeric[NT]) process(params RuleParams) ([]string, error) {
	schemaData := params.DataKey.(DataObject)
	var err error
	var bags []string

	schema := params.Schema
	originalData := params.OriginalData
	key := params.Key

	switch reflect.TypeOf(schema) {
	case reflect.TypeOf(Numeric[int]{}):
		if scMap, ok := schema.(Numeric[int]); ok {
			bags, err = scMap.validate(originalData, schemaData[key], params)
		}
	case reflect.TypeOf(Numeric[int32]{}):
		if scMap, ok := schema.(Numeric[int32]); ok {
			bags, err = scMap.validate(originalData, schemaData[key], params)
		}
	case reflect.TypeOf(Numeric[int64]{}):
		if scMap, ok := schema.(Numeric[int64]); ok {
			bags, err = scMap.validate(originalData, schemaData[key], params)
		}
	case reflect.TypeOf(Numeric[uint]{}):
		if scMap, ok := schema.(Numeric[uint]); ok {
			bags, err = scMap.validate(originalData, schemaData[key], params)
		}
	case reflect.TypeOf(Numeric[uint32]{}):
		if scMap, ok := schema.(Numeric[uint32]); ok {
			bags, err = scMap.validate(originalData, schemaData[key], params)
		}
	case reflect.TypeOf(Numeric[uint64]{}):
		if scMap, ok := schema.(Numeric[uint64]); ok {
			bags, err = scMap.validate(originalData, schemaData[key], params)
		}
	case reflect.TypeOf(Numeric[float32]{}):
		if scMap, ok := schema.(Numeric[float32]); ok {
			bags, err = scMap.validate(originalData, schemaData[key], params)
		}
	case reflect.TypeOf(Numeric[float64]{}):
		if scMap, ok := schema.(Numeric[float64]); ok {
			bags, err = scMap.validate(originalData, schemaData[key], params)
		}
	}
	return bags, err
	// pathKey := params.PathKey + key
	// if err != nil {
	// 	params.ErrorBags.append(pathKey, bags)
	// 	if options.AbortEarly {
	// 		return errors.New("error")
	// 	}
	// }
	// return nil
}

func (s Numeric[NT]) validate(jsonSource []byte, value any, params RuleParams) ([]string, error) {
	var bags []string
	key := params.Key
	option := params.Option

	err := s.assertRequired(key, value, &bags)

	if err != nil {
		return bags, err
	}

	err = s.assertRequiredIf(jsonSource, key, value, &bags)

	if err != nil {
		return bags, err
	}

	err = s.assertRequiredUnless(jsonSource, key, value, &bags)

	if err != nil {
		return bags, err
	}

	if value != nil {

		parsedValue, err := s.assertType(key, value, &bags)

		if err != nil {
			return bags, err
		}

		if err := s.assertMin(key, parsedValue, &bags); option.AbortEarly && err != nil {
			return bags, err
		}

		if err := s.assertMax(key, parsedValue, &bags); option.AbortEarly && err != nil {
			return bags, err
		}

		if err := s.assertMinDigits(key, parsedValue, &bags); option.AbortEarly && err != nil {
			return bags, err
		}

		if err := s.assertMaxDigits(key, parsedValue, &bags); option.AbortEarly && err != nil {
			return bags, err
		}

		if err := s.assertRegex(key, parsedValue, &bags); option.AbortEarly && err != nil {
			return bags, err
		}

		if err := s.assertNotRegex(key, parsedValue, &bags); option.AbortEarly && err != nil {
			return bags, err
		}

		if err := s.assertIn(key, parsedValue, &bags); option.AbortEarly && err != nil {
			return bags, err
		}

		if err := s.assertNotIn(key, parsedValue, &bags); option.AbortEarly && err != nil {
			return bags, err
		}

		if s.Custom != nil {
			if err := s.assertCustomValidation(s.Custom, jsonSource, parsedValue, params, &bags); option.AbortEarly && err != nil {
				return bags, err
			}
		}

	}

	if len(bags) > 0 {
		return bags, NumericValidationError
	}

	return bags, nil
}

func (s Numeric[NT]) assertType(key string, value any, bags *[]string) (NT, error) {
	if numericValue, ok := value.(NT); ok {
		return numericValue, nil
	}
	appendErrorBags(
		bags,
		fmt.Sprintf("%s must be type of %T", key, *new(NT)),
		"",
	)
	return 0, NumericValidationError
}

func (s Numeric[NT]) assertRequired(key string, value any, bags *[]string) error {
	if s.Required {
		if value == nil {
			appendErrorBags(
				bags,
				fmt.Sprintf("%s is required", key),
				s.Message.Required,
			)
			return NumericValidationError
		}
		if parsedValue, ok := isNumericValue[NT](value); ok {
			if digitLength(parsedValue) == 0 {
				appendErrorBags(
					bags,
					fmt.Sprintf("%s is required", key),
					s.Message.Required,
				)
				return NumericValidationError
			}
		}
	}
	return nil
}

func (s Numeric[NT]) assertRequiredIf(jsonSource []byte, key string, value any, bags *[]string) error {
	if s.RequiredIf != nil {
		if value == nil {
			comparedValue := gjson.GetBytes(jsonSource, s.RequiredIf.FieldPath)
			if comparedValue.Value() == s.RequiredIf.Value {
				appendErrorBags(
					bags,
					fmt.Sprintf("%s is required", key),
					s.Message.RequiredIf,
				)
				return NumericValidationError
			}
		} else {
			if parsedValue, ok := isNumericValue[NT](value); ok && digitLength(parsedValue) == 0 {
				comparedValue := gjson.GetBytes(jsonSource, s.RequiredIf.FieldPath)
				if comparedValue.Value() == s.RequiredIf.Value {
					appendErrorBags(
						bags,
						fmt.Sprintf("%s is required", key),
						s.Message.RequiredIf,
					)
					return NumericValidationError
				}
			}
		}
	}
	return nil
}

func (s Numeric[NT]) assertRequiredUnless(jsonSource []byte, key string, value any, bags *[]string) error {
	if s.RequiredUnless != nil {
		if value == nil {
			comparedValue := gjson.GetBytes(jsonSource, s.RequiredUnless.FieldPath)
			if comparedValue.Value() != s.RequiredUnless.Value {
				appendErrorBags(
					bags,
					fmt.Sprintf("%s is required", key),
					s.Message.RequiredUnless,
				)
				return NumericValidationError
			}
		} else {
			if parsedValue, ok := isNumericValue[NT](value); ok && digitLength(parsedValue) == 0 {
				comparedValue := gjson.GetBytes(jsonSource, s.RequiredUnless.FieldPath)
				if comparedValue.Value() != s.RequiredUnless.Value {
					appendErrorBags(
						bags,
						fmt.Sprintf("%s is required", key),
						s.Message.RequiredUnless,
					)
					return NumericValidationError
				}
			}
		}
	}
	return nil
}

func (s Numeric[NT]) assertMin(key string, value NT, bags *[]string) error {
	if s.Min > 0 && value < NT(s.Min) {
		appendErrorBags(
			bags,
			fmt.Sprintf("%s must be minimum of %d", key, s.Min),
			s.Message.Min,
		)
		return NumericValidationError
	}
	return nil
}

func (s Numeric[NT]) assertMax(key string, value NT, bags *[]string) error {
	if s.Max > 0 && value > NT(s.Max) {
		appendErrorBags(
			bags,
			fmt.Sprintf("%s must be maximum of %d", key, s.Max),
			s.Message.Max,
		)
		return NumericValidationError
	}
	return nil
}

func (s Numeric[NT]) assertMinDigits(key string, value NT, bags *[]string) error {
	if s.MinDigits > 0 && NT(digitLength(value)) < NT(s.MinDigits) {
		appendErrorBags(
			bags,
			fmt.Sprintf("%s total digits must be minimum of %d digit(s)", key, s.MinDigits),
			s.Message.MinDigits,
		)
		return NumericValidationError
	}
	return nil
}

func (s Numeric[NT]) assertMaxDigits(key string, value NT, bags *[]string) error {
	if s.MaxDigits > 0 && NT(digitLength(value)) > NT(s.MaxDigits) {
		appendErrorBags(
			bags,
			fmt.Sprintf("%s total digits must be maximum of %d digit(s)", key, s.MaxDigits),
			s.Message.MaxDigits,
		)
		return NumericValidationError
	}
	return nil
}

func (s Numeric[NT]) assertRegex(key string, value NT, bags *[]string) error {
	regx, err := regexp.Compile(s.Regex)
	if s.Regex != "" && (err != nil || !regx.MatchString(numericToString(value))) {
		appendErrorBags(
			bags,
			fmt.Sprintf("%s is not a valid format", key),
			s.Message.Regex,
		)
		return NumericValidationError
	}
	return nil
}

func (s Numeric[NT]) assertNotRegex(key string, value NT, bags *[]string) error {
	regx, err := regexp.Compile(s.Regex)
	if s.NotRegex != "" && digitLength(value) > 0 && (err != nil || regx.MatchString(numericToString(value))) {
		appendErrorBags(
			bags,
			fmt.Sprintf("%s is not a valid format", key),
			s.Message.NotRegex,
		)
		return NumericValidationError
	}
	return nil
}

func (s Numeric[NT]) assertIn(key string, value NT, bags *[]string) error {
	if len(s.In) > 0 && digitLength(value) > 0 && !slices.Contains(s.In, value) {
		var stringIn []string
		for _, n := range s.In {
			stringIn = append(stringIn, numericToString(n))
		}
		appendErrorBags(
			bags,
			fmt.Sprintf("%s must in %s", key, strings.Join(stringIn, ", ")),
			s.Message.In,
		)
		return NumericValidationError
	}
	return nil
}

func (s Numeric[NT]) assertNotIn(key string, value NT, bags *[]string) error {
	if len(s.In) > 0 && digitLength(value) > 0 && slices.Contains(s.In, value) {
		var stringNotIn []string
		for _, n := range s.NotIn {
			stringNotIn = append(stringNotIn, numericToString(n))
		}
		appendErrorBags(
			bags,
			fmt.Sprintf("%s must not in %s", key, strings.Join(stringNotIn, ", ")),
			s.Message.NotIn,
		)
		return NumericValidationError
	}
	return nil
}

func (s Numeric[NT]) assertCustomValidation(fc func(v NT, params RuleParams, look Lookup) error, jsonSource []byte, value NT, params RuleParams, bags *[]string) error {
	err := fc(value, params, func(k string) gjson.Result {
		return gjson.GetBytes(jsonSource, k)
	})
	if err != nil {
		appendErrorBags(
			bags,
			err.Error(),
			s.Message.Custom,
		)
		return NumericValidationError
	}
	return nil
}

func checkSchema() {

}

func isNumericValue[N NumericValue](value any) (N, bool) {
	if parsedValue, ok := value.(N); ok {
		return parsedValue, true
	}
	return 0, false
}

func numericToString[N NumericValue](value N) (stringValue string) {
	switch reflect.TypeOf(value).String() {
	case reflect.Int.String(), reflect.Int32.String(), reflect.Int64.String():
		stringValue = strconv.Itoa(int(value))
	case reflect.Float32.String():
		stringValue = strconv.FormatFloat(float64(value), 'g', -1, 32)
	case reflect.Float64.String():
		stringValue = strconv.FormatFloat(float64(value), 'g', -1, 64)
	}
	return
}

func digitLength[N NumericValue](value N) int {
	var stringValue = numericToString(value)
	return len([]rune(strings.Replace(stringValue, ".", "", -1)))
}
