package forms

import (
	log "github.com/sirupsen/logrus"
	"time"
)

// represents a form and the senders to send items in that form to
type Form struct {
	Name         string
	Destinations []Destination
	Validators   []Validator
}

// send a submission to all the Destinations in a form
func (form Form) SendSubmission(submission Submission) (Submission, error) {
	if submission.MessageIDS == nil {
		submission.MessageIDS = make(map[string]string)
	}

	err := form.ValidateSubmission(submission)
	if err != nil {
		return submission, err
	}

	for _, dest := range form.Destinations {
		logger := log.WithField("submission", submission).WithField("destination", dest)

		if messageid, ok := submission.MessageIDS[dest.Name]; ok {
			logger.Info("messageid already existed editing")
			err = dest.Edit(messageid, submission)
			if err != nil {
				logger.WithError(err).Error("EditSubmission returned an error")
			}
			continue
		}

		messageid, err := dest.Send(submission)
		if err != nil {
			logger.WithError(err).Error("SendSubmission returned an error")
			continue
		}

		submission.MessageIDS[dest.Name] = messageid
	}
	return submission, nil
}

// send a submission to all the Destinations in a form
func (form Form) DeleteSubmission(submission Submission) error {
	if submission.MessageIDS == nil {
		submission.MessageIDS = make(map[string]string)
	}

	for _, dest := range form.Destinations {
		logger := log.WithField("submission", submission).WithField("destination", dest)
		messageid := submission.MessageIDS[dest.Name]
		if messageid == "" {
			continue
		}

		err := dest.Delete(messageid)
		if err != nil {
			logger.WithError(err).Error("dest.Delete returned an error")
			continue
		}
	}

	submission.Deleted = true
	submission.DeletedAt = time.Now()
	return nil
}

// check if a Submission is valid
func (form *Form) ValidateSubmission(submission Submission) error {
	for _, validator := range form.Validators {
		err := validator.Validate(submission)
		if err != nil {
			return err
		}
	}

	return nil
}
