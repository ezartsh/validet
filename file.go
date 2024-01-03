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
}

type File struct {
	Required       bool
	RequiredIf     *RequiredIf
	RequiredUnless *RequiredUnless
	Max            int64
	Min            int64
	Mimes          string
	Message        FileErrorMessage
}

func (s *File) validate(jsonSource string, key string, value any, option Options) ([]string, error) {
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

		if err := s.assertMin(key, parsedValue, &bags); option.AbortEarly && err != nil {
			return bags, err
		}

		if err := s.assertMax(key, parsedValue, &bags); option.AbortEarly && err != nil {
			return bags, err
		}
	}

	if len(bags) > 0 {
		return bags, FileValidationError
	}

	return bags, nil

}

func (s *File) assertRequired(key string, value any, bags *[]string) error {
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

func (s *File) assertType(key string, value any, bags *[]string) (multipart.FileHeader, error) {
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

func (s *File) assertRequiredIf(jsonSource string, key string, value any, bags *[]string) error {
	if s.RequiredIf != nil && value == nil {
		comparedValue := gjson.Get(jsonSource, s.RequiredIf.FieldPath)
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

func (s *File) assertRequiredUnless(jsonSource string, key string, value any, bags *[]string) error {
	if s.RequiredUnless != nil && value == nil {
		comparedValue := gjson.Get(jsonSource, s.RequiredUnless.FieldPath)
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

func (s *File) assertMin(key string, value multipart.FileHeader, bags *[]string) error {
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

func (s *File) assertMax(key string, value multipart.FileHeader, bags *[]string) error {
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
