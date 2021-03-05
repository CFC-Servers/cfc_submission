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
	err := table.Get("UUID", identifier).Filter("'Deleted'<>?", true).One(&submission)
	return submission, err
}

func GetOwnerSubmissions(table dynamo.Table, ownerId string) ([]forms.Submission, error) {
	var submissions []forms.Submission
	// TODO pass this index in an environment variable
	err := table.Get("OwnerID", ownerId).
		Filter("'Deleted'<>?", true).
		Index("ownerid-createdat-index").
		Order(dynamo.Descending).
		All(&submissions)
	return submissions, err
}
