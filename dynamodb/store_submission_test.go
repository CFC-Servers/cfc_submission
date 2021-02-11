package dynamodb

import (
	"flag"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/cfc-servers/cfc_suggestions/forms"
	"github.com/guregu/dynamo"
	"github.com/stretchr/testify/assert"
	"testing"
)

var runDynamoDBTests = flag.Bool("usedynamo", false, "run tests using dynamodb")

var exampleSubmission = forms.Submission{
	FormName: "suggestion",
	UUID:     "560dff48-8bad-4606-9b3d-dcce902214cf",
	OwnerInfo: forms.OwnerInfo{
		Name:   "HMM#0001",
		ID:     "179237013373845504",
		Avatar: "https://cdn.discordapp.com/avatars/179237013373845504/cef7ff8bf178f6a1d9552cc68a4e0620.png",
	},
	MessageIDS: nil,
	Fields:     nil,
}

func TestPutSubmission(t *testing.T) {
	if !*runDynamoDBTests {
		t.Skip()
	}

	s, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	assert.NoError(t, err)

	db := dynamo.New(s)
	table := db.Table("cfc_forms_submissions")

	err = PutSubmission(table, exampleSubmission)
	assert.NoError(t, err)
}

func TestGetSubmission(t *testing.T) {
	if !*runDynamoDBTests {
		t.Skip()
	}

	s, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	assert.NoError(t, err)

	db := dynamo.New(s)
	table := db.Table("cfc_forms_submissions")

	err = PutSubmission(table, exampleSubmission)
	assert.NoError(t, err)

	submission, err := GetSubmission(table, exampleSubmission.UUID)
	assert.NoError(t, err)

	assert.Equal(t, submission.OwnerInfo.Avatar, exampleSubmission.OwnerInfo.Avatar)

}
