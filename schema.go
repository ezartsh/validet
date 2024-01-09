package validet

type RuleParams struct {
	OriginalData []byte
	DataKey      any
	PathKey      []string
	Key          string
	Schema       Rule
	ErrorBags    *ErrorBag
	Option       Options
}

type Rule interface {
	validate(source []byte, value any, params RuleParams) ([]string, error)
	isMyTypeOf(schema any) bool
	process(RuleParams) ([]string, error)
}

type SchemaRules = map[string]Rule

type schemaContainer struct {
	Data    DataObject
	Items   SchemaRules
	Options Options
}

func NewSchema(d DataObject, items SchemaRules, options Options) schemaContainer {
	return schemaContainer{
		Data:    d,
		Items:   items,
		Options: options,
	}
}

func (s *schemaContainer) Validate() (ErrorBag, error) {
	return validate(s.Data, s.Items, s.Options)
}
