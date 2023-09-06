package validet

import (
	"encoding/json"
	"errors"
	"slices"
	"testing"
)

func Test_String_Required(t *testing.T) {
	t.Run("it should error when the property is not exist", func(t *testing.T) {
		schema := String{Required: true}
		_, err := schema.validate("", "test", "", Options{})
		if !errors.Is(err, StringValidationError) {
			t.Errorf("Actual = %v, Expected = %v", err, StringValidationError)
		}
	})
	t.Run("it should error when the property value is null", func(t *testing.T) {
		schema := String{Required: true}
		_, err := schema.validate("", "test", nil, Options{})
		if !errors.Is(err, StringValidationError) {
			t.Errorf("Actual = %v, Expected = %v", err, StringValidationError)
		}
	})
	t.Run("it should error when the property value is empty string", func(t *testing.T) {
		schema := String{Required: true}
		_, err := schema.validate("", "test", "", Options{})
		if !errors.Is(err, StringValidationError) {
			t.Errorf("Actual = %v, Expected = %v", err, StringValidationError)
		}
	})
}

func Test_String_Value_Type(t *testing.T) {
	t.Run("it should error when the property value is not a string", func(t *testing.T) {
		schema := String{Required: true}
		_, err := schema.validate("", "test", 123, Options{})
		if !errors.Is(err, StringValidationError) {
			t.Errorf("Actual = %v, Expected = %v", err, StringValidationError)
		}
	})
}

func Test_String_RequiredIf(t *testing.T) {
	t.Run("it should error when the property value is empty and another property value is x", func(t *testing.T) {
		data := DataObject{
			"test_1": "x",
			"test":   "",
		}
		jsonString, _ := json.Marshal(data)
		schema := String{RequiredIf: &RequiredIf{
			FieldPath: "test_1",
			Value:     "x",
		}}
		_, err := schema.validate(string(jsonString), "test", "", Options{})
		if !errors.Is(err, StringValidationError) {
			t.Errorf("Actual = %v, Expected = %v", err, StringValidationError)
		}
	})
	t.Run("it should not error when the property value is empty and another property value is not x", func(t *testing.T) {
		data := DataObject{
			"test_1": "y",
			"test":   "",
		}
		jsonString, _ := json.Marshal(data)
		schema := String{RequiredIf: &RequiredIf{
			FieldPath: "test_1",
			Value:     "x",
		}}
		_, err := schema.validate(string(jsonString), "test", "", Options{})
		if errors.Is(err, StringValidationError) {
			t.Errorf("Actual = %v, Expected = %v", err, nil)
		}
	})
}

func Test_String_RequiredUnless(t *testing.T) {
	t.Run("it should error when the property value is empty and another property value is not x", func(t *testing.T) {
		data := DataObject{
			"test_1": "y",
			"test":   "",
		}
		jsonString, _ := json.Marshal(data)
		schema := String{RequiredUnless: &RequiredUnless{
			FieldPath: "test_1",
			Value:     "x",
		}}
		_, err := schema.validate(string(jsonString), "test", "", Options{})
		if !errors.Is(err, StringValidationError) {
			t.Errorf("Actual = %v, Expected = %v", err, StringValidationError)
		}
	})
	t.Run("it should not error when the property value is empty and another property value is x", func(t *testing.T) {
		data := DataObject{
			"test_1": "x",
			"test":   "",
		}
		jsonString, _ := json.Marshal(data)
		schema := String{RequiredUnless: &RequiredUnless{
			FieldPath: "test_1",
			Value:     "x",
		}}
		_, err := schema.validate(string(jsonString), "test", "", Options{})
		if errors.Is(err, StringValidationError) {
			t.Errorf("Actual = %v, Expected = %v", err, nil)
		}
	})
}

func Test_String_Min(t *testing.T) {
	t.Run("it should error when the length of property value is not bigger than or equal to x", func(t *testing.T) {
		schema := String{Min: 2}
		_, err := schema.validate("", "test", "x", Options{})
		if !errors.Is(err, StringValidationError) {
			t.Errorf("Actual = %v, Expected = %v", err, StringValidationError)
		}
	})
}

func Test_String_Max(t *testing.T) {
	t.Run("it should error when the length of property value is not less than or equal to x", func(t *testing.T) {
		schema := String{Max: 2}
		_, err := schema.validate("", "test", "xxxx", Options{})
		if !errors.Is(err, StringValidationError) {
			t.Errorf("Actual = %v, Expected = %v", err, StringValidationError)
		}
	})
}

func Test_String_Regex(t *testing.T) {
	t.Run("it should error when the property value is not match the expression", func(t *testing.T) {
		schema := String{Regex: "^test$"}
		_, err := schema.validate("", "test", "test2", Options{})
		if !errors.Is(err, StringValidationError) {
			t.Errorf("Actual = %v, Expected = %v", err, StringValidationError)
		}
	})
}

func Test_String_NotRegex(t *testing.T) {
	t.Run("it should error when the property value is match the expression", func(t *testing.T) {
		schema := String{NotRegex: "^test$"}
		_, err := schema.validate("", "test", "test", Options{})
		if !errors.Is(err, StringValidationError) {
			t.Errorf("Actual = %v, Expected = %v", err, StringValidationError)
		}
	})
}

func Test_String_In(t *testing.T) {
	t.Run("it should error when the property value is not on the list", func(t *testing.T) {
		schema := String{In: []string{"one", "two", "three"}}
		_, err := schema.validate("", "test", "four", Options{})
		if !errors.Is(err, StringValidationError) {
			t.Errorf("Actual = %v, Expected = %v", err, StringValidationError)
		}
	})
}

func Test_String_NotIn(t *testing.T) {
	t.Run("it should error when the property value is on the list", func(t *testing.T) {
		schema := String{NotIn: []string{"one", "two", "three"}}
		_, err := schema.validate("", "test", "one", Options{})
		if !errors.Is(err, StringValidationError) {
			t.Errorf("Actual = %v, Expected = %v", err, StringValidationError)
		}
	})
}

func Test_String_Email(t *testing.T) {
	cases := []string{
		"test",
		"test@test",
		"test.com",
		"test@test@test.com",
	}
	for _, cs := range cases {
		t.Run("it should error when the property value is not valid email e.g "+cs, func(t *testing.T) {
			schema := String{Email: true}
			_, err := schema.validate("", "test", cs, Options{})
			if !errors.Is(err, StringValidationError) {
				t.Errorf("Actual = %v, Expected = %v", err, StringValidationError)
			}
		})
	}
}

func Test_String_Alpha(t *testing.T) {
	t.Run("it should error when the property value is not alphabetical", func(t *testing.T) {
		schema := String{Alpha: true}
		_, err := schema.validate("", "test", "test123", Options{})
		if !errors.Is(err, StringValidationError) {
			t.Errorf("Actual = %v, Expected = %v", err, StringValidationError)
		}
	})
}

func Test_String_AlphaNumeric(t *testing.T) {
	t.Run("it should error when the property value is not alphabetical and numeric", func(t *testing.T) {
		schema := String{AlphaNumeric: true}
		_, err := schema.validate("", "test", "_test_", Options{})
		if !errors.Is(err, StringValidationError) {
			t.Errorf("Actual = %v, Expected = %v", err, StringValidationError)
		}
	})
}

func Test_String_Url(t *testing.T) {
	cases := []string{
		"test.com",
		"http//test.com",
		"https:test.com",
		"http://test",
	}
	for _, cs := range cases {
		t.Run("it should error when the property value is not valid url e.g "+cs, func(t *testing.T) {
			schema := String{Url: &Url{Http: true, Https: true}}
			_, err := schema.validate("", "test", cs, Options{})
			if !errors.Is(err, StringValidationError) {
				t.Errorf("Actual = %v, Expected = %v", err, StringValidationError)
			}
		})
	}
}

func Test_String_Custom_Validation(t *testing.T) {
	t.Run("it should error when the custom validation return error", func(t *testing.T) {
		schema := String{Custom: func(v string) error {
			return StringValidationError
		}}
		_, err := schema.validate("", "test", "_test_", Options{})
		if !errors.Is(err, StringValidationError) {
			t.Errorf("Actual = %v, Expected = %v", err, StringValidationError)
		}
	})
}

func Test_String_Custom_Message(t *testing.T) {
	t.Run("it should return custom message error when the custom message is configured.", func(t *testing.T) {
		schema := String{Min: 5, Message: StringErrorMessage{Min: "minimum 2"}}
		bags, _ := schema.validate("", "test", "tst", Options{})
		if !slices.Contains(bags, "minimum 2") {
			t.Fatalf("Actual = %v, Expected contain = minimum 2", bags)
		}
	})
}
