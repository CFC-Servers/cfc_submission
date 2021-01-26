package actions

import (
	"github.com/cfc-servers/cfc_suggestions/forms"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateSubmission(t *testing.T) {
	if !*runDynamoDBTests { t.Skip() }
	table := getTable(t)

	sub, err := CreateSubmission(table, "ban-appeal", forms.Owner{
		Name:       "FluffieFoxBoi#9583",
		Identifier: "345010836089339906",
		Avatar:     "https://api.foxorsomething.net/fox/d3e60f8a-9622-44ef-a0e3-16507293ed3c.png",
	})
	assert.NoError(t, err)

	assert.NotEqual(t, sub.UUID, "")
}

