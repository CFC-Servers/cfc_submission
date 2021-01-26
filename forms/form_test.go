package forms

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockSender struct {
	mock.Mock
}

func (m MockSender) Send(submission Submission) (string, error) {
	args := m.Called(submission)
	return args.String(0), args.Error(1)
}

var TestSubmissions = []struct{
	submission        Submission
	expectedMessageId string
	expectedError     error
}{
	{
		Submission{
			FormName: "suggestion",
			UUID:     "7184a4dc-c6e7-418f-a61f-8eeb7801c11e",
			Owner: Owner{
				Name:       "HMM#001",
				Identifier: "433348813462175745",
				Avatar:     "https://cdn.discordapp.com/avatars/179237013373845504/cef7ff8bf178f6a1d9552cc68a4e0620.png",
			},
			MessageIDS: map[string]string{},
			Fields: SubmissionFields{
				"_description": "description1",
				"_title":       "title2",
			},
		}, "12345", nil,
	},
	{
		Submission{
			FormName: "suggestion",
			UUID:     "7184a4dc-c6e7-418f-a61f-8eeb7801c11e",
			Owner: Owner{
				Name:       "HMM#001",
				Identifier: "433348813462175745",
				Avatar:     "https://cdn.discordapp.com/avatars/179237013373845504/cef7ff8bf178f6a1d9552cc68a4e0620.png",
			},
			MessageIDS: map[string]string{},
			Fields: SubmissionFields{
				"_description": "description1",
				"_title":       "title2",
			},
		}, "12345", errors.New("test error"),
	},
}



func TestForm_SendSubmission(t *testing.T) {
	sender := MockSender{}

	testDestination := Destination{
		Name:  "test",
		Sender: &sender,
	}

	form := Form{
		Name: "suggestion",
		Destinations: []Destination{testDestination},
	}

	for _, data := range TestSubmissions {
		sender.On("Send", data.submission).Return(data.expectedMessageId, data.expectedError)

		submission, _ := form.SendSubmission(data.submission)
		if data.expectedError == nil {
			assert.Contains(t, submission.MessageIDS, testDestination.Name)
		} else {
			assert.NotContains(t, submission.MessageIDS, testDestination.Name)
		}

	}
}