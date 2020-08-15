package storage

type Suggestion struct {
	Identifier  string `json:"identifier"`
	Owner       string `json:"owner"`
	Active      bool   `json:"active"`
	ContentJson string `json:"content_json"`
}

type SuggestionStore interface {
	Create(owner string) (*Suggestion, error)
	Get(identifier string) (*Suggestion, error)
	Delete(identifier string) error
	Update(identifier string, active bool, contentJson string) error
	DeleteActive(owner string) error
}
