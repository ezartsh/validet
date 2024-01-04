package validet

import "errors"

type ErrorBag struct {
	Errors map[string][]string
	Status bool
}

func NewErrorBags() *ErrorBag {
	return &ErrorBag{
		Errors: make(map[string][]string),
	}
}

func (e *ErrorBag) add(key string, m string) {
	e.Errors[key] = append(e.Errors[key], m)
}

func (e *ErrorBag) append(key string, msgs []string) {
	if mv, ok := e.Errors[key]; ok {
		e.Errors[key] = append(mv, msgs...)
	} else {
		e.Errors[key] = msgs
	}
}

var ObjectValidationError = errors.New("object validation failed")
var SliceObjectValidationError = errors.New("slice object validation failed")
var StringValidationError = errors.New("string validation failed")
var NumericValidationError = errors.New("numeric validation failed")
var SliceValidationError = errors.New("slice validation failed")
var FileValidationError = errors.New("file validation failed")
var BooleanValidationError = errors.New("boolean validation failed")

var ErrorRequiredField = errors.New("field cannot be empty")
