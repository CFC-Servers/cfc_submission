package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cfc-servers/cfc_suggestions/app"
	"github.com/cfc-servers/cfc_suggestions/dynamodb"
	"github.com/cfc-servers/cfc_suggestions/forms"
	"github.com/cfc-servers/cfc_suggestions/util"
	"github.com/guregu/dynamo"
	"net/http"
)

func main() {
	lambda.Start(SendSubmissionHandler)
}

func SendSubmissionHandler(req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	uuid := req.PathParameters["uuid"]

	submission, err := dynamodb.GetSubmission(util.GetTable(), uuid)
	if err != nil {
		return ErrorResponse(err), nil
	}

	submission.Fields = make(forms.SubmissionFields)
	if err := json.Unmarshal([]byte(req.Body), &submission.Fields); err != nil {
		return util.Response(http.StatusBadRequest, fmt.Sprintf(`{"Error": "%v"}`, err)), err
	}

	form, _ := app.GetForm(submission.FormName)

	submission, err = form.SendSubmission(submission)
	if err != nil {
		return ErrorResponse(err), nil
	}
	err = dynamodb.PutSubmission(util.GetTable(), submission)
	
	return util.Response(http.StatusOK, submission), err
}

func ErrorResponse(err error) events.APIGatewayV2HTTPResponse {
	if errors.Is(err, dynamo.ErrNotFound) {
		return util.Response(http.StatusNotFound, map[string]string{"Error": "not found"})
	}
	if errors.Is(err, forms.ValidationErr) {
		return util.Response(http.StatusBadRequest, map[string]string{"Error": err.Error()})
	}
	return util.Response(http.StatusInternalServerError, err)
}