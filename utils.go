package validet

import (
	"fmt"
	"reflect"
)

func appendErrorBags(bags *[]string, om string, cm string) {
	errorBags := *bags
	message := om
	if cm != "" {
		message = cm
	}
	*bags = append(errorBags, message)
}

func isObjectValue(value any) bool {
	return reflect.TypeOf(value) == reflect.TypeOf(DataObject{})
}

func isStringValue(value any) bool {
	return reflect.TypeOf(value) == reflect.TypeOf("")
}

func isBooelanValue(value any) bool {
	return reflect.TypeOf(value) == reflect.TypeOf(true)
}

func stringLength(value any) int {
	if value == nil {
		return 0
	}
	stringValue := value.(string)
	return len([]rune(stringValue))
}

func msgf(format string, v ...any) string {
	return fmt.Sprintf(format, v...)
}
