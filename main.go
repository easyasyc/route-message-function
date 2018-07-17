package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

type Message struct {
	Message string `json:"message"`
}

func HandleRequest(ctx context.Context, message Message) (string, error) {
	return fmt.Sprintf("Hello %s!", message.Message), nil
}

func main() {
	lambda.Start(HandleRequest)
}
