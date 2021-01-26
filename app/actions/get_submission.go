package actions

import (
	"github.com/cfc-servers/cfc_suggestions/dynamodb"
	"github.com/cfc-servers/cfc_suggestions/forms"
	"github.com/guregu/dynamo"
)

func GetSubmission(table dynamo.Table, identifier string) (submission forms.Submission, err error) {
	return dynamodb.GetSubmission(table, identifier)
}