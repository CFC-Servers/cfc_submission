package sqlite

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/cfc-servers/cfc_suggestions/suggestions"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
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
			content_json TEXT
            
		)`)
	if err != nil {
		panic(err)
	}
	return &SqliteSuggestionsStore{
		db: db,
	}
}

func (store *SqliteSuggestionsStore) Create(owner string) (*suggestions.Suggestion, error) {
	suggestion := suggestions.Suggestion{
		Identifier: newIdentifier(),
		Owner:      owner,
		Sent:       false,
	}

	_, err := store.db.Exec(
		"INSERT INTO cfc_suggestions(identifier, owner) VALUES(?, ?)",
		suggestion.Identifier, suggestion.Owner, false)

	if err != nil {
		return nil, err
	}

	return &suggestion, nil
}

func (store *SqliteSuggestionsStore) Delete(owner string, onlyUnsent bool) error {
	query := "DELETE FROM cfc_suggestions WHERE owner=?"
	if onlyUnsent == true {
		query = query + " AND sent=0"
	}
	_, err := store.db.Exec(query, owner)
	return err
}

func (store *SqliteSuggestionsStore) Get(identifier string) (*suggestions.Suggestion, error) {
	suggestion := suggestions.Suggestion{}
	row := store.db.QueryRow("SELECT * FROM cfc_suggestions WHERE identifier=?", identifier)

	var contentJson []byte
	var sentInt int
	err := row.Scan(&suggestion.Identifier, &suggestion.Owner, &sentInt, &suggestion.MessageID, &contentJson)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, err
	} else if err != nil {
		log.Errorf("Database error: %v", err)
	}
	suggestion.Sent = sentInt != 0

	json.Unmarshal(contentJson, &suggestion.Content)
	return &suggestion, err
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
