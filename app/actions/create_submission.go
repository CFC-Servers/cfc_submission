package actions

import (
	"github.com/cfc-servers/cfc_suggestions/dynamodb"
	"github.com/cfc-servers/cfc_suggestions/forms"
	"github.com/google/uuid"
	"github.com/guregu/dynamo"
)

// Creates and saves a new submission
func CreateSubmission(table dynamo.Table, formGetter func(string) (forms.Form, error), FormName string, owner forms.Owner) (forms.Submission, error){
	_,  err := formGetter(FormName)
	if err != nil {
		return forms.Submission{}, err
	}
	submission := forms.Submission{
		FormName:  FormName,
		Owner: owner,
		UUID:    uuid.New().String(),
	}

	err = dynamodb.PutSubmission(table, submission)

	return submission, err




}