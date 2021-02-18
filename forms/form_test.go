package forms

import (
	"errors"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
)

var testForm = Form{

	Name: "suggestion",
	Validators: []Validator{
		FieldValidator{}.Accept(
			Field("shortField").Min(1).Max(10),
			Field("longField").Min(10).Max(50),
		),
	},

	Destinations: []Destination{
		Destination{"mock", fakeSender},
	},

	Formatter: DefaultFormatter{Color: 0x34eb5b},
}
var fakeSender = &MockDestination{}

func TestForm_SendSubmission(t *testing.T) {
	fakeSender.On("Send", mock.Anything).Return("1234", nil)

	testData := map[string]struct {
		ExpectedError error
		Fields        SubmissionFields
	}{
		"always_pass": {
			ExpectedError: nil,
			Fields: SubmissionFields{
				"shortField": "short",
				"longField":  "this field must be very long",
			},
		},
		"fail_on_shortField": {
			ExpectedError: ErrFieldTooLong,
			Fields: SubmissionFields{
				"shortField": "this is toooo long",
				"longField":  "this field must be very long",
			},
		},
		"fail_on_longField": {
			ExpectedError: ErrFieldTooShort,
			Fields: SubmissionFields{
				"shortField": "short",
				"longField":  "too short",
			},
		},
		"fail_on_missing_field": {
			ExpectedError: ErrMissingField,
			Fields: SubmissionFields{
				"shortField": "short",
			},
		},
		"fail_on_unkown_field": {
			ExpectedError: ErrUnknownField,
			Fields: SubmissionFields{
				"shortField": "this is toooo long",
				"longField":  "this field must be very long",
				"weirdField": "test",
			},
		},
	}

	for k, submissionData := range testData {
		t.Run(k, func(t *testing.T) {
			submission := NewSubmission(testForm, OwnerInfo{
				ID:   "1234",
				Name: "HMM",
			})
			submission.Fields = submissionData.Fields
			submission, err := testForm.SendSubmission(submission)
			if !errors.Is(err, submissionData.ExpectedError) {
				t.Errorf("Unexpected err response from SendSubmission %v", err)
			}

			if submissionData.ExpectedError != nil {
				return
			}

			if reflect.DeepEqual(submission.Content, FormattedContent{}) {
				t.Error("empty content on non error submission")
			}
		})
	}
}

type MockDestination struct {
	mock.Mock
}

func (dest *MockDestination) Send(submission Submission) (string, error) {
	args := dest.Called(submission)
	return args.String(0), args.Error(1)
}

func (dest *MockDestination) Edit(messageid string, submission Submission) error {
	args := dest.Called(messageid, submission)
	return args.Error(0)
}

func (dest *MockDestination) Delete(messageid string) error {
	args := dest.Called(messageid)
	return args.Error(0)
}
