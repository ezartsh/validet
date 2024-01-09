package validet

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/tidwall/gjson"
)

type SliceObjectErrorMessage struct {
	Required       string
	RequiredIf     string
	RequiredUnless string
	Min            string
	Max            string
	Custom         string
}

type SliceObject struct {
	Required       bool
	RequiredIf     *RequiredIf
	RequiredUnless *RequiredUnless
	Min            int
	Max            int
	Item           DataObject
	Custom         func(v []DataObject, look Lookup) error
	Message        SliceObjectErrorMessage
}

func (s SliceObject) isMyTypeOf(schema any) bool {
	return reflect.TypeOf(schema).Kind() == reflect.Struct && reflect.TypeOf(schema) == reflect.TypeOf(SliceObject{})
}

func (s SliceObject) process(params RuleParams) ([]string, error) {
	schemaData := params.DataKey.(DataObject)
	// var err error
	// var bags []string

	errorBags := *params.ErrorBags
	schema := params.Schema
	originalData := params.OriginalData
	key := params.Key
	options := params.Option

	if scSliceObject, ok := schema.(SliceObject); ok {
		bags, err := scSliceObject.validate(originalData, key, schemaData[key], options)
		if err != nil {
			return bags, err
			// errorBags.append(params.PathKey+key, bags)
			// if options.AbortEarly {
			// 	return errors.New("new error")
			// }
		} else {
			schemaDataValues := schemaData[key].([]interface{})
			for i, value := range schemaDataValues {
				for scObjItemKey, scObjItemValue := range scSliceObject.Item {
					mapSchemas(originalData, params.PathKey+key+"."+strconv.Itoa(i), scObjItemKey, value, scObjItemValue, &errorBags, options)
					// if options.AbortEarly && len(errorBags.Errors) > 0 {
					// 	return errors.New("new error")
					// }
				}
			}
		}
	}

	// pathKey := params.PathKey + key
	// if err != nil {
	// 	params.ErrorBags.append(pathKey, bags)
	// 	if options.AbortEarly {
	// 		return errors.New("error")
	// 	}
	// }
	return []string{}, nil
}

func (s SliceObject) validate(jsonSource []byte, key string, value any, option Options) ([]string, error) {
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

		parsedValue, err := s.assertType(key, value, &bags)

		if err != nil {
			return bags, err
		}

		if len(parsedValue) > 0 {

			if err != nil {
				return bags, err
			}

			if err := s.assertMin(key, parsedValue, &bags); option.AbortEarly && err != nil {
				return bags, err
			}

			if err := s.assertMax(key, parsedValue, &bags); option.AbortEarly && err != nil {
				return bags, err
			}

			if s.Custom != nil {
				if err := s.assertCustomValidation(s.Custom, jsonSource, parsedValue, &bags); option.AbortEarly && err != nil {
					return bags, err
				}
			}

		}

	}

	if len(bags) > 0 {
		return bags, SliceObjectValidationError
	}

	return bags, nil

}

func (s SliceObject) assertRequired(key string, value any, bags *[]string) error {
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

func (s SliceObject) assertType(key string, value any, bags *[]string) ([]DataObject, error) {
	if values, ok := value.([]interface{}); ok {
		sliceDataObject := []DataObject{}
		for _, v := range values {
			if value, ok := v.(DataObject); ok {
				sliceDataObject = append(sliceDataObject, value)
			}
		}
		if len(values) == len(sliceDataObject) {
			return sliceDataObject, nil
		}
	}
	appendErrorBags(
		bags,
		fmt.Sprintf("%s must be type of data object", key),
		"",
	)
	return []DataObject{}, SliceObjectValidationError
}

func (s SliceObject) assertRequiredIf(jsonSource []byte, key string, value any, bags *[]string) error {
	values := value.([]interface{})
	if s.RequiredIf != nil && (value == nil || len(values) == 0) {
		comparedValue := gjson.GetBytes(jsonSource, s.RequiredIf.FieldPath)
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

func (s SliceObject) assertRequiredUnless(jsonSource []byte, key string, value any, bags *[]string) error {
	values := value.([]interface{})
	if s.RequiredUnless != nil && (value == nil || len(values) == 0) {
		comparedValue := gjson.GetBytes(jsonSource, s.RequiredUnless.FieldPath)
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

func (s SliceObject) assertMin(key string, values []DataObject, bags *[]string) error {
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

func (s SliceObject) assertMax(key string, values []DataObject, bags *[]string) error {
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

func (s SliceObject) assertCustomValidation(fc func(v []DataObject, look Lookup) error, jsonSource []byte, value []DataObject, bags *[]string) error {
	err := fc(value, func(k string) gjson.Result {
		return gjson.GetBytes(jsonSource, k)
	})
	if err != nil {
		appendErrorBags(
			bags,
			err.Error(),
			s.Message.Custom,
		)
		return SliceObjectValidationError
	}
	return nil
}
