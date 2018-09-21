package document

import "bytes"

// ValidationErrors store each errors associated to every fields of a model
type ValidationErrors struct {
	Errors map[string][]string `json:"errors"`
}

func (v *ValidationErrors) Error() string {
	var buffer bytes.Buffer

	for field, errors := range v.Errors {
		buffer.WriteString(field)
		buffer.WriteString("=")
		for _, err := range errors {
			buffer.WriteString(err)
			buffer.WriteString(", ")
		}
	}
	return buffer.String()
}

// ValidationErrorsBuilder is used to provide a simple way to create a ValidationErrors struct. The typical usecase is:
//	func (m *MyModel) Validate(ctx context.Context) *ValidationErrors {
//		validations := document.NewValidationErrorsBuilder()
//
//		if m.Name == "" {
//			validations.Set("name", "should not be empty")
//		}
//
//		if m.Email == "" {
//			validations.Set("email", "should not be empty")
//		}
//
//		return validations.Build()
//	}
type ValidationErrorsBuilder struct {
	errors map[string][]string
}

// NewValidationErrors return an empty ValidationErrors struct
func NewValidationErrorsBuilder() *ValidationErrorsBuilder {
	return &ValidationErrorsBuilder{
		errors: make(map[string][]string),
	}
}

// Set will add an error on a specific field, if the field already contains an error, it will just add it to the current errors list
func (v *ValidationErrorsBuilder) Set(field, err string) {
	v.errors[field] = append(v.errors[field], err)
}

// Get will return all errors set for a specific field
func (v *ValidationErrorsBuilder) Get(field string) []string {
	return v.errors[field]
}

// Build will send a ValidationErrors struct if there is some errors or nil if no errors has been defined
func (v *ValidationErrorsBuilder) Build() *ValidationErrors {
	if len(v.errors) == 0 {
		return nil
	}

	return &ValidationErrors{
		Errors: v.errors,
	}
}
