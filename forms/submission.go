package forms

import (
	"github.com/google/uuid"
	"time"
)

// submission
type Submission struct {
	UUID      string
	FormName  string
	OwnerID   string
	OwnerInfo OwnerInfo
	CreatedAt time.Time

	MessageIDS map[string]string
	Fields     SubmissionFields

	Deleted   bool
	DeletedAt time.Time
}

type OwnerInfo struct {
	ID     string
	Name   string
	Avatar string
}

func NewSubmission(form Form, ownerInfo OwnerInfo) Submission {
	return Submission{
		UUID:      uuid.New().String(),
		FormName:  form.Name,
		OwnerID:   ownerInfo.ID,
		OwnerInfo: ownerInfo,
		CreatedAt: time.Now(),
	}
}

type SubmissionFields map[string]interface{}

func (fields SubmissionFields) Has(key string) bool {
	_, ok := fields[key]
	return ok
}

func (fields SubmissionFields) Get(key string) string {
	value, ok := fields[key]
	if !ok {
		return ""
	}

	if strValue, ok := value.(string); ok {
		return strValue
	}

	return ""
}
