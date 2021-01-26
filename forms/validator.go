package forms

import (
	"errors"
	"fmt"
)

type Validator interface {
	Validate(Submission) (Submission, error)
}

var ValidationErr = errors.New("error validating request")
var MissingFieldErr = fmt.Errorf("%w: missing required field", ValidationErr)
var FieldTooShortErr  = fmt.Errorf("%w: field too short", ValidationErr)
var FieldTooLongErr = fmt.Errorf("%w: field too long", ValidationErr)

type FieldValidator struct {
	Fields []FieldValidatorField
}

type FieldValidatorField struct {
	Name string
	MinLength int
	MaxLength int
	IsOptional bool
	ValidOptions map[string]bool
}

func(f FieldValidator) Accept(field ...FieldValidatorField) FieldValidator {
	if f.Fields == nil {
		f.Fields = make([]FieldValidatorField, 0)
	}
	f.Fields = append(f.Fields, field...)
	return f
}

func (f FieldValidator) Validate(submission Submission) (Submission, error) {
	for _, field := range f.Fields {
		if !submission.Fields.Has(field.Name) && !field.IsOptional {
			return submission, fmt.Errorf("%w: field %v is required", MissingFieldErr, field.Name)
		}


		strField := submission.Fields.Get(field.Name)
		if field.MinLength != 0 && len(strField) < field.MinLength {
			return submission, fmt.Errorf("%w: field must be shorter than %v characters", ValidationErr, field.MinLength)
		}

		if field.MaxLength != 0 && len(strField) > field.MaxLength {
			return submission,  fmt.Errorf("%w: field must be longer than %v characters", ValidationErr, field.MaxLength)
		}

	}

	return submission, nil
}

func Field(name string) FieldValidatorField {
	return FieldValidatorField{
		Name: name,
	}
}

func (field FieldValidatorField) Min(min int) FieldValidatorField{
	field.MinLength = min
	return field
}

func  (field FieldValidatorField) Max(max int) FieldValidatorField {
	field.MaxLength = max
	return field
}

func (field FieldValidatorField) MustBe(options ...string) FieldValidatorField {
	if field.ValidOptions  == nil  {
		field.ValidOptions  = make(map[string]bool)
	}
	for _, option := range options {
		field.ValidOptions[option] = true
	}
	return field
}
func (field FieldValidatorField) Optional() FieldValidatorField {
	field.IsOptional = true
	return field
}