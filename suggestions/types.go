package suggestions

import "errors"

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
	Delete(identifier string) error
	Update(suggestion *Suggestion) error
	DeleteByOwner(owner string, onlyUnsent bool) error
}

var ErrMessageNotFound = errors.New("Message not found")

type Destination interface {
	Send(suggestion *Suggestion) (messageId string, err error)
	SendEdit(suggestion *Suggestion) (messageId string, err error)
	Delete(messageId string) error
}
