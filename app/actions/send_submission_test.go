package actions

import (
	"github.com/cfc-servers/cfc_suggestions/app"
	"github.com/cfc-servers/cfc_suggestions/forms"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSendSubmission(t *testing.T) {
	if !*runDynamoDBTests { t.Skip() }

	table := getTable(t)
	submission, err := SendSubmission(table, app.GetForm, "37e2b0b9-42e4-4d65-964c-83cac8258821", forms.SubmissionFields{
		"why":         "i like addon aaaaaaaaaaaaaaaa",
		"whyNot":      "i dont like addon aaaaaaaaaaaaaaaaaaaaaaa",
		"description": "is a very cool addon aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, submission.MessageIDS, "messageids table is empty")
}