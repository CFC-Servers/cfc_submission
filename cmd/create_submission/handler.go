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
	"net/http"
)

func main() {
	lambda.Start(CreateSubmissionHandler)
}

func CreateSubmissionHandler(req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	var data struct {
		FormName string
		Owner    forms.OwnerInfo
	}
	if err := json.Unmarshal([]byte(req.Body), &data); err != nil {
		return util.Response(http.StatusBadRequest, fmt.Sprintf(`{"Error": "%v"}`, err)), err
	}

	form, err := app.GetForm(data.FormName)
	if err != nil {
		return ErrorResponse(err), nil
	}

	submission := forms.NewSubmission(form, data.Owner)

	err = dynamodb.PutSubmission(util.GetTable(), submission)
	if err != nil {
		panic(err)
	}

	return util.Response(http.StatusCreated, submission), nil
}

func ErrorResponse(err error) events.APIGatewayV2HTTPResponse {
	if errors.Is(err, app.ErrMissingForm) {
		return util.Response(http.StatusBadRequest, fmt.Sprintf(`{"Error": "Invalid form name"}`))
	}

	return util.Response(http.StatusInternalServerError, err)
}
