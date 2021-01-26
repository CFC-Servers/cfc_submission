package forms

// submission
type Submission struct {
	FormName string
	UUID string
	Owner Owner
	MessageIDS map[string]string
	Fields SubmissionFields
}

type Owner struct {
	Name string
	Identifier string
	Avatar string
}

type SubmissionFields map[string]interface{}

func (fields SubmissionFields) Has(key string)  bool {
	_, ok := fields[key]
	return ok
}
func (fields SubmissionFields) Get(key string) string {
	value, ok := fields[key]
	if !ok {
		return ""
	}


	if strValue, ok := value.(string); ok {
		return strValue
	}

	return ""
}


