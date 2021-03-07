package util

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"os"
)

func GetTable() dynamo.Table {
	s, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"), // TODO this  region shouldnt be hardcoded
	})
	if err != nil {
		panic(err)
	}

	db := dynamo.New(s)
	return db.Table(os.Getenv("SUBMISSIONS_DATABASE"))
}

func Response(StatusCode int, obj interface{}) events.APIGatewayV2HTTPResponse {
	var body string
	switch obj := obj.(type) {
	case string:
		body = obj
	case byte:
		body = string(obj)
	case error:
		data, err := json.Marshal(map[string]interface{}{
			"Error": obj.Error(),
		})
		if err != nil {
			panic(err)
		}

		body = string(data)
	default:
		data, err := json.Marshal(obj)
		if err != nil {
			panic(err)
		}

		body = string(data)
	}
	return events.APIGatewayV2HTTPResponse{
		StatusCode: StatusCode,
		Body:       body,
		Headers:    map[string]string{"Content-Type": "application/json"},
	}
}
