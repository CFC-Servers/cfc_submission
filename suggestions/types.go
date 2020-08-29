package suggestions

import (
	"errors"
	"time"
)

type Suggestion struct {
	Identifier string             `json:"identifier"`
	Owner      string             `json:"owner"`
	Sent       bool               `json:"sent"`
	MessageID  string             `json:"message_id"`
	Content    *SuggestionContent `json:"content"`
	CreatedAt  time.Time          `json:"created_at"`
}

type SuggestionContent struct {
	Realm     string `json:"realm"`
	Link      string `json:"link" length:"0-124"`
	Title     string `json:"title" length:"6-124"`
	Why       string `json:"why" length:"18-1024"`
	WhyNot    string `json:"whyNot" length:"18-1024"`
	Anonymous bool   `json:"anonymous"`
}

type SuggestionStore interface {
	Create(suggestion *Suggestion) (*Suggestion, error)
	GetWhere(map[string]interface{}) ([]*Suggestion, error)
	DeleteWhere(map[string]interface{}) error
	Update(suggestion *Suggestion) error
}

var ErrMessageNotFound = errors.New("Message not found")

type Destination interface {
	Send(suggestion *Suggestion) (messageId string, err error)
	SendEdit(suggestion *Suggestion) (messageId string, err error)
	Delete(messageId string) error
}
