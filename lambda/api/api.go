package api

import (
	"encoding/json"
	"fmt"
	"lambda-func/database"
	"lambda-func/types"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type ApiHandler struct {
	dbStore database.UserStore
}

func NewApiHandler(dbStore database.UserStore) ApiHandler {
	return ApiHandler{
		dbStore: dbStore,
	}
}

func (api ApiHandler) RegisterUserHandler(
	req events.APIGatewayProxyRequest,
) (events.APIGatewayProxyResponse, error) {
	var registerUser types.RegisterUser

	err := json.Unmarshal([]byte(req.Body), &registerUser)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Invalid request",
			StatusCode: http.StatusBadRequest,
		}, err
	}

	if registerUser.Username == "" || registerUser.Password == "" {
		return events.APIGatewayProxyResponse{
			Body:       "Invalid request",
			StatusCode: http.StatusBadRequest,
		}, err
	}

	userExist, err := api.dbStore.DoesUserExist(registerUser.Username)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal server error",
			StatusCode: http.StatusBadRequest,
		}, err
	}

	if userExist {
		return events.APIGatewayProxyResponse{
			Body:       "User already exist",
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	user, err := types.NewUser(registerUser)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Bad request",
			StatusCode: http.StatusBadRequest,
		}, err
	}

	err = api.dbStore.UserInset(user)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal server error",
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	return events.APIGatewayProxyResponse{
		Body:       "Register success",
		StatusCode: http.StatusOK,
	}, err
}

func (api ApiHandler) LoginUserHandler(
	req events.APIGatewayProxyRequest,
) (events.APIGatewayProxyResponse, error) {

	var loginUser types.LoginUser

	err := json.Unmarshal([]byte(req.Body), &loginUser)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Invalid request",
			StatusCode: http.StatusBadRequest,
		}, err
	}

	if loginUser.Username == "" || loginUser.Password == "" {
		return events.APIGatewayProxyResponse{
			Body:       "Invalid request",
			StatusCode: http.StatusBadRequest,
		}, err
	}

	user, err := api.dbStore.GetUser(loginUser.Username)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal server error",
			StatusCode: http.StatusBadRequest,
		}, err
	}

	match := types.ValidatePassword(user.Password, loginUser.Password)

	if !match {
		return events.APIGatewayProxyResponse{
			Body:       "Wrong login or password",
			StatusCode: http.StatusUnauthorized,
		}, err
	}

	accessToken := types.CreateToken(user.Username)
	successMsg := fmt.Sprintf(`{"access_token": "%s"}`, accessToken)

	return events.APIGatewayProxyResponse{
		Body:       successMsg,
		StatusCode: http.StatusOK,
	}, err
}

func (api ApiHandler) ProtestedRoute(
	req events.APIGatewayProxyRequest,
) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body:       "This is protected route",
		StatusCode: http.StatusOK,
	}, nil
}
