package main

import (
	"fmt"
	"lambda-func/app"

	"github.com/aws/aws-lambda-go/lambda"
)

type MyEvent struct {
	UserName string `json:"username"`
}

func HandleRequest(event MyEvent) (string, error) {
	if event.UserName == "" {
		return "", fmt.Errorf("username can not be empty")
	}

	return fmt.Sprintf("Successfully called by: %s", event.UserName), nil
}

func main() {
	myApp := app.NewApp()
	lambda.Start(myApp.ApiHandler.RegisterUserHandler)
}
