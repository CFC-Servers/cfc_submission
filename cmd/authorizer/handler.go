package main

import (
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(AuthorizerRequest)
}

var IsValidToken = make(map[string]bool)

func init() {
	validTokens := strings.Split(os.Getenv("VALID_AUTH_TOKENS"), ",")

	for _, token := range validTokens {
		if strings.TrimSpace(token) == "" {
			continue
		}
		IsValidToken[token] = true
	}
}

type AuthResponse struct {
	IsAuthorized bool `json:"isAuthorized"`
}

var AllowedResponse = AuthResponse{IsAuthorized: true}
var DeniedResponse = AuthResponse{IsAuthorized: false}

func AuthorizerRequest(req events.APIGatewayCustomAuthorizerRequest) (AuthResponse, error) {
	if IsValidToken[req.AuthorizationToken] {
		return AllowedResponse, nil
	}

	return DeniedResponse, nil
}
