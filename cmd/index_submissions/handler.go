package main

import (
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cfc-servers/cfc_suggestions/dynamodb"
	"github.com/cfc-servers/cfc_suggestions/util"
	"github.com/guregu/dynamo"
	"net/http"
)

func main() {
	lambda.Start(IndexSubmissionsHandler)
}

func IndexSubmissionsHandler(req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	ownerId := req.QueryStringParameters["owner"]
	if ownerId == "" {
		return util.Response(http.StatusBadRequest, "an owner must be provided"), nil
	}

	submission, err := dynamodb.GetOwnerSubmissions(util.GetTable(), ownerId)
	if err != nil {
		return errorResponse(err), nil
	}

	return util.Response(http.StatusOK, submission), err
}

func errorResponse(err error) events.APIGatewayV2HTTPResponse {
	if errors.Is(err, dynamo.ErrNotFound) {
		return util.Response(http.StatusNotFound, map[string]string{"Error": "not found"})
	}

	return util.Response(http.StatusInternalServerError, err)
}
