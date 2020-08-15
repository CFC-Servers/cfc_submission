package storage

type Suggestion struct {
	Identifier string `json:"identifier"`
	Owner      string `json:"owner"`
}

type SuggestionStore interface {
	Create(owner string) (*Suggestion, error)
	Get(identifier string) (*Suggestion, error)
	Delete(identifier string) error
}
