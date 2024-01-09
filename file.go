package validet

import (
	"fmt"
	"mime/multipart"

	"github.com/tidwall/gjson"
)

type FileErrorMessage struct {
	Required       string
	RequiredIf     string
	RequiredUnless string
	Max            string
	Min            string
	Mimes          string
	Custom         string
}

type File struct {
	Required       bool
	RequiredIf     *RequiredIf
	RequiredUnless *RequiredUnless
	Max            int64
	Min            int64
	Mimes          string
	Custom         func(v multipart.FileHeader, path PathKey, look Lookup) error
	Message        FileErrorMessage
}

func (s File) validate(source []byte, value any, params RuleParams) ([]string, error) {
	var bags []string
	key := params.Key
	option := params.Option

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

		if s.Custom != nil {
			if err := s.assertCustomValidation(s.Custom, source, parsedValue, PathKey{
				Previous: params.PathKey,
				Current:  params.Key,
			}, &bags); option.AbortEarly && err != nil {
				return bags, err
			}
		}
	}

	if len(bags) > 0 {
		return bags, FileValidationError
	}

	return bags, nil

}

func (s File) assertRequired(key string, value any, bags *[]string) error {
	if s.Required {
		if value == nil {
			appendErrorBags(
				bags,
				fmt.Sprintf("%s is required", key),
				s.Message.Required,
			)
			return FileValidationError
		}
	}
	return nil
}

func (s File) assertType(key string, value any, bags *[]string) (multipart.FileHeader, error) {
	if parsedValue, ok := value.(multipart.FileHeader); ok {
		return parsedValue, nil
	}

	appendErrorBags(
		bags,
		fmt.Sprintf("%s must be a type of file", key),
		s.Message.Required,
	)
	return multipart.FileHeader{}, FileValidationError
}

func (s File) assertRequiredIf(jsonSource []byte, key string, value any, bags *[]string) error {
	if s.RequiredIf != nil && value == nil {
		comparedValue := gjson.GetBytes(jsonSource, s.RequiredIf.FieldPath)
		if comparedValue.String() == s.RequiredIf.Value {
			appendErrorBags(
				bags,
				fmt.Sprintf("%s is required", key),
				s.Message.RequiredIf,
			)
			return FileValidationError
		}
	}
	return nil
}

func (s File) assertRequiredUnless(jsonSource []byte, key string, value any, bags *[]string) error {
	if s.RequiredUnless != nil && value == nil {
		comparedValue := gjson.GetBytes(jsonSource, s.RequiredUnless.FieldPath)
		if comparedValue.String() != s.RequiredUnless.Value {
			appendErrorBags(
				bags,
				fmt.Sprintf("%s is required", key),
				s.Message.RequiredUnless,
			)
			return FileValidationError
		}
	}
	return nil
}

func (s File) assertMin(key string, value multipart.FileHeader, bags *[]string) error {
	if value.Size < s.Min {
		appendErrorBags(
			bags,
			fmt.Sprintf("%s size must be at minimum %d", key, s.Min),
			s.Message.Min,
		)
		return FileValidationError
	}
	return nil
}

func (s File) assertMax(key string, value multipart.FileHeader, bags *[]string) error {
	if value.Size > s.Max {
		appendErrorBags(
			bags,
			fmt.Sprintf("%s size must be at maximum %d", key, s.Max),
			s.Message.Max,
		)
		return FileValidationError
	}
	return nil
}

func (s File) assertCustomValidation(fc func(v multipart.FileHeader, path PathKey, look Lookup) error, jsonSource []byte, value multipart.FileHeader, path PathKey, bags *[]string) error {
	err := fc(value, path, func(k string) gjson.Result {
		return gjson.GetBytes(jsonSource, k)
	})
	if err != nil {
		appendErrorBags(
			bags,
			err.Error(),
			s.Message.Custom,
		)
		return FileValidationError
	}
	return nil
}
