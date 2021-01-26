package actions

import (
	"flag"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"github.com/stretchr/testify/assert"
	"testing"
)

var runDynamoDBTests = flag.Bool("usedynamo", false, "run tests using dynamodb")

func getTable(t *testing.T) dynamo.Table {
	s, err := session.NewSession( &aws.Config{
		Region: aws.String("us-east-1"),
	})
	assert.NoError(t, err)
	db := dynamo.New(s)
	table := db.Table("cfc_forms_submissions")
	return table
}
