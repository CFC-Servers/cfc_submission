package sqlite

import (
	"database/sql"
	"encoding/json"
	"github.com/cfc-servers/cfc_suggestions/suggestions"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

type SqliteSuggestionsStore struct {
	db         *sql.DB
	LogQueries bool
}

func NewStore(file string) *SqliteSuggestionsStore {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS cfc_suggestions(
			identifier TEXT NOT NULL PRIMARY KEY, 
			owner TEXT NOT NULL,
			sent SMALLINT NOT NULL DEFAULT 0,
			message_id TEXT DEFAULT '',
			content_json TEXT,
			created_at INTEGER DEFAULT 0
        )
	`)
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
	suggestion.CreatedAt = time.Now()

	_, err := store.exec(
		"INSERT INTO cfc_suggestions(identifier, owner, created_at) VALUES(?, ?, ?)",
		suggestion.Identifier, suggestion.Owner, suggestion.CreatedAt.Unix())

	if err != nil {
		return nil, err
	}

	return suggestion, nil
}

func (store *SqliteSuggestionsStore) DeleteWhere(conditions map[string]interface{}) error {
	where, values := constructWhere(conditions)
	query := "DELETE FROM cfc_suggestions" + where

	_, err := store.exec(query, values...)
	return err
}

func (store *SqliteSuggestionsStore) GetWhere(conditions map[string]interface{}) ([]*suggestions.Suggestion, error) {
	outputSuggestions := make([]*suggestions.Suggestion, 0)

	after, _ := conditions["after"].(int64)
	delete(conditions, "after")

	where, values := constructWhere(conditions, "created_at>?")
	values = append(values, after)

	query := "SELECT * FROM cfc_suggestions" + where

	rows, err := store.query(query, values...)
	if err != nil {
		return outputSuggestions, err
	}

	for rows.Next() {
		outputSuggestions = append(outputSuggestions, scanIntoSuggestion(rows))
	}

	return outputSuggestions, nil
}

func (store *SqliteSuggestionsStore) Update(suggestion *suggestions.Suggestion) error {
	contentJson, _ := json.Marshal(suggestion.Content)
	sentInt := 0
	if suggestion.Sent {
		sentInt = 1
	}

	_, err := store.exec(
		"UPDATE cfc_suggestions SET content_json=?, sent=?, message_id=? WHERE identifier=?",
		contentJson, sentInt, suggestion.MessageID, suggestion.Identifier,
	)

	return err
}

func (store *SqliteSuggestionsStore) query(query string, args ...interface{}) (*sql.Rows, error) {
	if store.LogQueries {
		log.Infof("Sqlite query: \"%v\"  %v", query, args)
	}
	out, err := store.db.Query(query, args...)
	if err != nil {
		log.Infof("Query errored: %v", err)
	}
	return out, err
}

func (store *SqliteSuggestionsStore) exec(query string, args ...interface{}) (sql.Result, error) {
	if store.LogQueries {
		log.Infof("Sqlite query: \"%v\"  %v", query, args)
	}
	out, err := store.db.Exec(query, args...)
	if err != nil {
		log.Infof("Query errored: %v", err)
	}
	return out, err
}

func scanIntoSuggestion(rows *sql.Rows) *suggestions.Suggestion {
	suggestion := suggestions.Suggestion{}
	var contentJson []byte
	var sentInt int
	var createdAt int64
	rows.Scan(&suggestion.Identifier, &suggestion.Owner, &sentInt, &suggestion.MessageID, &contentJson, createdAt)
	suggestion.CreatedAt = time.Unix(createdAt, 0)
	if sentInt == 1 {
		suggestion.Sent = true
	}
	json.Unmarshal(contentJson, &suggestion.Content)
	return &suggestion
}

func newIdentifier() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

func constructWhere(conditions map[string]interface{}, extraConditions ...string) (string, []interface{}) {
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

	for _, condition := range extraConditions {
		if firstCondition {
			firstCondition = false
			queryBuilder.WriteString(" WHERE ")
		} else {
			queryBuilder.WriteString(" AND ")
		}
		queryBuilder.WriteString(condition)
	}

	return queryBuilder.String(), values
}
