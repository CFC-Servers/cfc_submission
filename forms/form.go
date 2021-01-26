package forms

import log "github.com/sirupsen/logrus"

// represents a form and the senders to send items in that form to
type Form struct {
	Name string
	Destinations []Destination
	Validators []Validator
}

// send a submission to all the Destinations in a form
func (form Form) SendSubmission(submission Submission) (Submission, error) {
	if submission.MessageIDS == nil {
		submission.MessageIDS = make(map[string]string)
	}

	var err error
	for _, validator := range form.Validators {
		submission, err = validator.Validate(submission)
		if err != nil {
			return submission, err
		}
	}

	for _, dest := range form.Destinations {
		logger := log.WithField("submission", submission).WithField("destination", dest)

		if _, ok := submission.MessageIDS[dest.Name]; ok {
			logger.Info("messageid already existed not sending")

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