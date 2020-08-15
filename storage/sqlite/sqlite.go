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

	db.Exec(`CREATE TABLE IF NOT EXISTS cfc_suggestions(identifier TEXT NOT NULL PRIMARY KEY, owner TEXT NOT NULL)`)

	return &SqliteSuggestionsStore{
		db: db,
	}
}

func (store *SqliteSuggestionsStore) Create(owner string) (*storage.Suggestion, error) {
	suggestion := storage.Suggestion{
		Identifier: newIdentifier(),
		Owner:      owner,
	}

	_, err := store.db.Exec("INSERT INTO cfc_suggestions(identifier, owner) VALUES(?, ?)", suggestion.Identifier, suggestion.Owner)
	if err != nil {
		return nil, err
	}

	return &suggestion, nil
}

func (store *SqliteSuggestionsStore) Get(identifier string) (*storage.Suggestion, error) {
	suggestion := storage.Suggestion{}
	row := store.db.QueryRow("SELECT identifier, owner FROM cfc_suggestions WHERE identifier=?", identifier)
	row.Scan(&suggestion.Identifier, &suggestion.Owner)
	return &suggestion, nil
}

func (store *SqliteSuggestionsStore) Delete(identifier string) error {
	_, err := store.db.Exec("DELETE FROM cfc_suggestions WHERE identifier=?", identifier)
	return err
}

func newIdentifier() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}
