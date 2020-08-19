package suggestions

type Suggestion struct {
	Identifier string             `json:"identifier"`
	Owner      string             `json:"owner"`
	Sent       bool               `json:"sent"`
	MessageID  string             `json:"message_id"`
	Content    *SuggestionContent `json:"content"`
	CreatedAt  string             `json:"created_at"`
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
	Create(suggestion *Suggestion) (*Suggestion, error)
	GetWhere(map[string]interface{}) ([]*Suggestion, error)
	DeleteWhere(map[string]interface{}) error
	Update(suggestion *Suggestion) error
}

type Destination interface {
	Send(suggestion *Suggestion) (messageId string, err error)
	SendEdit(suggestion *Suggestion) (messageId string, err error)
	Delete(messageId string) error
}
