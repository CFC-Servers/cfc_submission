package forms

// a Destination to send submitted forms to
type Destination struct {
	Name string
	Sender
}

// the Sender interface should handle sending, editing, and deleting messages
type Sender interface {
	Send(Submission) (messageid string, err error)
	Edit(string, Submission) (err error)
}
