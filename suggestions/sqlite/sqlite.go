package sqlite

import (
	"database/sql"
	"encoding/json"
	"github.com/cfc-servers/cfc_suggestions/suggestions"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"strings"
)

type SqliteSuggestionsStore struct {
	db *sql.DB
}

func NewStore(file string) *SqliteSuggestionsStore {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS cfc_suggestions(
			identifier TEXT NOT NULL PRIMARY KEY, 
			owner TEXT NOT NULL,
			sent SMALLINT NOT NULL DEFAULT 0,
			message_id TEXT DEFAULT '',
			content_json TEXT,
			created_at NOT NULL DEFAULT current_timestamp
            
		)`)
	if err != nil {
		panic(err)
	}
	return &SqliteSuggestionsStore{
		db: db,
	}
}

func (store *SqliteSuggestionsStore) Create(suggestion *suggestions.Suggestion) (*suggestions.Suggestion, error) {
	if suggestion.Identifier == "" {
		suggestion.Identifier = newIdentifier()
	}

	_, err := store.db.Exec(
		"INSERT INTO cfc_suggestions(identifier, owner) VALUES(?, ?)",
		suggestion.Identifier, suggestion.Owner, false)

	if err != nil {
		return nil, err
	}

	return suggestion, nil
}

func (store *SqliteSuggestionsStore) DeleteWhere(conditions map[string]interface{}) error {
	where, values := constructWhere(conditions)
	query := "DELETE FROM cfc_suggestions" + where

	_, err := store.db.Exec(query, values...)
	return err
}

func (store *SqliteSuggestionsStore) GetWhere(conditions map[string]interface{}) ([]*suggestions.Suggestion, error) {
	where, values := constructWhere(conditions)
	query := "SELECT * FROM cfc_suggestions" + where

	rows, _ := store.db.Query(query, values...)
	// TODO do something with errors

	outputSuggestions := make([]*suggestions.Suggestion, 0)
	for rows.Next() {
		suggestion := suggestions.Suggestion{}
		var contentJson []byte
		var sentInt int

		rows.Scan(&suggestion.Identifier, &suggestion.Owner, &sentInt, &suggestion.MessageID, &contentJson, &suggestion.CreatedAt)
		if sentInt == 1 {
			suggestion.Sent = true
		}
		json.Unmarshal(contentJson, &suggestion.Content)
		outputSuggestions = append(outputSuggestions, &suggestion)

	}

	return outputSuggestions, nil
}

func (store *SqliteSuggestionsStore) Update(suggestion *suggestions.Suggestion) error {
	contentJson, _ := json.Marshal(suggestion.Content)
	sentInt := 0
	if suggestion.Sent {
		sentInt = 1
	}

	_, err := store.db.Exec(
		"UPDATE cfc_suggestions SET content_json=?, sent=?, message_id=? WHERE identifier=?",
		contentJson, sentInt, suggestion.MessageID, suggestion.Identifier,
	)

	return err
}

func newIdentifier() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

func constructWhere(conditions map[string]interface{}) (string, []interface{}) {
	var queryBuilder strings.Builder
	var values []interface{}
	firstCondition := true
	for column, value := range conditions {
		if valueBool, ok := value.(bool); ok {
			if valueBool {
				value = 1
			} else {
				value = 0
			}
		}

		if firstCondition {
			firstCondition = false
			queryBuilder.WriteString(" WHERE ")
		} else {
			queryBuilder.WriteString(" AND ")
		}
		queryBuilder.WriteString(column)
		queryBuilder.WriteString("=?")
		values = append(values, value)

	}
	return queryBuilder.String(), values
}
