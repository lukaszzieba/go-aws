package main

import (
	"lambda-func/app"
	"lambda-func/middleware"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	myApp := app.NewApp()
	lambda.Start(
		func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
			switch request.Path {

			case "/register":
				return myApp.ApiHandler.RegisterUserHandler(request)

			case "/login":
				return myApp.ApiHandler.LoginUserHandler(request)

			case "/protected":
				return middleware.ValidateJWTMiddleware(myApp.ApiHandler.ProtestedRoute)(request)

			default:
				return events.APIGatewayProxyResponse{
					Body:       "Not found",
					StatusCode: http.StatusNotFound,
				}, nil
			}
		},
	)
}
