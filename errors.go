package validet

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
