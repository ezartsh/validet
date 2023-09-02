package validet

func appendErrorBags(bags *[]string, om string, cm string) {
	errorBags := *bags
	message := om
	if cm != "" {
		message = cm
	}
	*bags = append(errorBags, message)
}
