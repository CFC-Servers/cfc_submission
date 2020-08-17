package suggestions

type Suggestion struct {
	Identifier string             `json:"identifier"`
	Owner      string             `json:"owner"`
	Sent       bool               `json:"sent"`
	Content    *SuggestionContent `json:"content"`
	MessageID  string             `json:"message_id"`
}

type SuggestionContent struct {
	Realm     string `json:"realm"`
	Link      string `json:"link"`
	Title     string `json:"title"`
	Why       string `json:"why"`
	WhyNot    string `json:"whyNot"`
	Anonymous bool   `json:"anonymous"`
}

type SuggestionStore interface {
	Create(owner string) (*Suggestion, error)
	Get(identifier string) (*Suggestion, error)
	Update(suggestion *Suggestion) error
	Delete(owner string, onlyUnsent bool) error
}

type Destination interface {
	Send(content *Suggestion) (messageId string, err error)
	SendEdit(content *Suggestion) (messageId string, err error)
}
