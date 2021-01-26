package lambda_submission_handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/cfc-servers/cfc_suggestions/app/actions"
	"github.com/cfc-servers/cfc_suggestions/forms"
	"github.com/guregu/dynamo"
	"net/http"
	"os"
)

func HandleRequest(req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	table, err := getTable()
	if err != nil {
		panic(err)
	}

	switch req.RouteKey{
	case "POST /submissions":
		var data struct {
			FormName string
			Owner forms.Owner
		}

		err := json.Unmarshal([]byte(req.Body), &data)
		if err != nil {
			return response(http.StatusBadRequest, fmt.Sprintf(`{"error": "%v"}`, err)), nil
		}

		submission, err := actions.CreateSubmission(table, data.FormName, data.Owner)
		if errors.Is(err, forms.ValidationErr) {
			return response(http.StatusBadRequest, fmt.Sprintf(`{"error": "%v"}`, err)), nil
		}

		return response(http.StatusCreated, submission), nil

	case "POST /submissions/{id}/send":

	}

	return response(http.StatusNotFound, map[string]string{
		"Message": "not found",
		"RoutKey": req.RouteKey,
	}), nil
}


func getTable() (dynamo.Table, error) {
	s, err := session.NewSession( &aws.Config{
		Region: aws.String("us-east-1"), // TODO this  region shouldnt be hardcoded
	})
	if  err != nil {
		return dynamo.Table{}, err
	}

	db := dynamo.New(s)
	table := db.Table(os.Getenv("TABLE_NAME"))
	return table, nil
}
func response(StatusCode int, obj interface{})  events.APIGatewayV2HTTPResponse {
	var body string
	switch obj := obj.(type) {
	case string:
		body = obj
	case byte:
		body = string(obj)
	default:
		data, err := json.Marshal(obj)
		if  err != nil {
			panic(err)
		}

		body = string(data)
	}
	return events.APIGatewayV2HTTPResponse{
		StatusCode:        StatusCode,
		Body:              body,
	}
}