package main

import (
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cfc-servers/cfc_suggestions/app"
	"github.com/cfc-servers/cfc_suggestions/dynamodb"
	"github.com/cfc-servers/cfc_suggestions/util"
	"github.com/guregu/dynamo"
	"net/http"
)

func main() {
	lambda.Start(GetSubmissionHandler)
}

func GetSubmissionHandler(req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	uuid := req.PathParameters["uuid"]

	submission, err := dynamodb.GetSubmission(util.GetTable(), uuid)
	if err != nil {
		return errorResponse(err), nil
	}

	form, err := app.GetForm(submission.FormName)
	if err != nil {
		return errorResponse(err), nil
	}

	err = form.DeleteSubmission(submission)
	if err != nil {
		return errorResponse(err), nil
	}

	return util.Response(http.StatusNoContent, ""), nil
}

func errorResponse(err error) events.APIGatewayV2HTTPResponse {
	if errors.Is(err, dynamo.ErrNotFound) {
		return util.Response(http.StatusNotFound, map[string]string{"Error": "not found"})
	}

	return util.Response(http.StatusInternalServerError, err)
}
