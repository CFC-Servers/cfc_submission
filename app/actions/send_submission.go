package actions

import (
	"github.com/cfc-servers/cfc_suggestions/dynamodb"
	"github.com/cfc-servers/cfc_suggestions/forms"
	"github.com/guregu/dynamo"
)



// Sends a submission to its  destinations
func SendSubmission(table dynamo.Table, formGetter func(string) (forms.Form, error), identifier string, values forms.SubmissionFields) (forms.Submission, error){
	submission, err := dynamodb.GetSubmission(table, identifier)
	if err != nil {
		return forms.Submission{}, err
	}

	form, err := formGetter(submission.FormName)
	if err != nil {
		return forms.Submission{}, err
	}

	submission.Fields = values

	submission, err = form.SendSubmission(submission)
	if err != nil {
		return submission, err
	}

	err = dynamodb.PutSubmission(table, submission)
	return submission, err
}