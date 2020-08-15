package sqlite

import (
	"database/sql"
	"github.com/cfc-servers/cfc_suggestions/storage"
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

	db.Exec(
		`CREATE TABLE IF NOT EXISTS cfc_suggestions(
			identifier TEXT NOT NULL PRIMARY KEY, 
			owner TEXT NOT NULL,
			active SMALLINT NOT NULL,
			content_json TEXT
		)`)

	return &SqliteSuggestionsStore{
		db: db,
	}
}

func (store *SqliteSuggestionsStore) Create(owner string) (*storage.Suggestion, error) {
	suggestion := storage.Suggestion{
		Identifier: newIdentifier(),
		Owner:      owner,
		Active:     true,
	}

	_, err := store.db.Exec("INSERT INTO cfc_suggestions(identifier, owner, active) VALUES(?, ?, 1)", suggestion.Identifier, suggestion.Owner)
	if err != nil {
		return nil, err
	}

	return &suggestion, nil
}

func (store *SqliteSuggestionsStore) DeleteActive(owner string) error {
	_, err := store.db.Exec("DELETE FROM cfc_suggestions WHERE owner=? AND active=1", owner)
	return err
}

func (store *SqliteSuggestionsStore) GetActive(identifier string) (*storage.Suggestion, error) {
	suggestion := storage.Suggestion{}
	row := store.db.QueryRow("SELECT * FROM cfc_suggestions WHERE identifier=? AND active=1", identifier)
	err := row.Scan(&suggestion.Identifier, &suggestion.Owner, &suggestion.Active, &suggestion.ContentJson)
	return &suggestion, err
}

func (store *SqliteSuggestionsStore) Delete(identifier string) error {
	_, err := store.db.Exec("DELETE FROM cfc_suggestions WHERE identifier=?", identifier)
	return err
}

func (store *SqliteSuggestionsStore) Update(identifier string, active bool, contentJson string) error {
	activeInt := 0
	if active {
		activeInt = 1
	}

	_, err := store.db.Exec(
		"UPDATE cfc_suggestions SET content_json=?, active=? WHERE identifier=?",
		contentJson, activeInt, identifier,
	)

	return err
}

func newIdentifier() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}
