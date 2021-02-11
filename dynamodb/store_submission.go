package dynamodb

import (
	"github.com/cfc-servers/cfc_suggestions/forms"
	"github.com/guregu/dynamo"
)

func PutSubmission(table dynamo.Table, submission forms.Submission) error {
	return table.Put(submission).Run()
}

func GetSubmission(table dynamo.Table, identifier string) (forms.Submission, error) {
	var submission forms.Submission
	err := table.Get("UUID", identifier).One(&submission)
	return submission, err
}

func DeleteSubmission(table dynamo.Table, identifier string) error {
	return table.Delete("UUID", identifier).Run()
}
