package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"os"
	"strings"
)

func main() {
	lambda.Start(AuthorizerRequest)
}

var IsValidToken map[string]bool

func init() {
	validTokens := strings.Split(os.Getenv("VALID_AUTH_TOKENS"), ",")

	for _, token := range validTokens {
		IsValidToken[token] = true
	}
}

type AuthResponse struct {
	IsAuthorized bool `json:"isAuthorized"`
}

var AllowedResponse = AuthResponse{IsAuthorized: true}
var DeniedResponse = AuthResponse{IsAuthorized: false}

func AuthorizerRequest(req events.APIGatewayCustomAuthorizerRequest) AuthResponse {
	if IsValidToken[req.AuthorizationToken] {
		return AllowedResponse
	}

	return DeniedResponse
}
