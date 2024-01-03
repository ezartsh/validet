package validet

import "reflect"

func appendErrorBags(bags *[]string, om string, cm string) {
	errorBags := *bags
	message := om
	if cm != "" {
		message = cm
	}
	*bags = append(errorBags, message)
}

func isStringValue(value any) bool {
	return reflect.TypeOf(value) == reflect.TypeOf("")
}

func stringLength(value any) int {
	if value == nil {
		return 0
	}
	stringValue := value.(string)
	return len([]rune(stringValue))
}
