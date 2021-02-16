package forms

import (
	"errors"
	"fmt"
)

type Validator interface {
	Validate(Submission) error
}

var ValidationErr = errors.New("error validating request")
var MissingFieldErr = fmt.Errorf("%w: missing required field", ValidationErr)
var FieldTooShortErr = fmt.Errorf("%w: field too short", ValidationErr)
var FieldTooLongErr = fmt.Errorf("%w: field too long", ValidationErr)
var UnknownFieldErr = fmt.Errorf("%w: unkown field", ValidationErr)
type FieldValidator struct {
	Fields map[string]FieldValidatorField
}

type FieldValidatorField struct {
	Name         string
	MinLength    int
	MaxLength    int
	IsOptional   bool
	ValidOptions map[string]bool
}

func (f FieldValidator) Accept(newFields ...FieldValidatorField) FieldValidator {
	if f.Fields == nil {
		f.Fields = make(map[string]FieldValidatorField)
	}

	for _, newField :=  range newFields {
		f.Fields[newField.Name] = newField
	}

	return f
}

func (f FieldValidator) Validate(submission Submission) error {
	for k, _ := range submission.Fields {
		_, ok := f.Fields[k]
		if !ok {
			return fmt.Errorf("%w: did not expect field %v", UnknownFieldErr, k)
		}
	}

	for _, validator := range f.Fields {
		if !submission.Fields.Has(validator.Name) && !validator.IsOptional {
			return fmt.Errorf("%w: field %v is required", MissingFieldErr, validator.Name)
		}

		strField := submission.Fields.GetString(validator.Name)
		if validator.MinLength != 0 && len(strField) < validator.MinLength {
			return fmt.Errorf("%w: field %v must be longer than %v characters", FieldTooShortErr, validator.Name, validator.MinLength)
		}

		if validator.MaxLength != 0 && len(strField) > validator.MaxLength {
			return fmt.Errorf("%w: field %v must be shorter than %v characters", FieldTooLongErr, validator.Name, validator.MaxLength)
		}
	}


	return nil
}

func Field(name string) FieldValidatorField {
	return FieldValidatorField{
		Name: name,
	}
}

func (field FieldValidatorField) Min(min int) FieldValidatorField {
	field.MinLength = min
	return field
}

func (field FieldValidatorField) Max(max int) FieldValidatorField {
	field.MaxLength = max
	return field
}

func (field FieldValidatorField) MustBe(options ...string) FieldValidatorField {
	if field.ValidOptions == nil {
		field.ValidOptions = make(map[string]bool)
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
