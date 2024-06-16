package api

import (
	"fmt"
	"lambda-func/database"
	"lambda-func/types"
)

type ApiHandler struct {
	dbStore database.DynamoDBClient
}

func NewApiHandler(dbStore database.DynamoDBClient) ApiHandler {
	return ApiHandler{
		dbStore: dbStore,
	}
}

func (api ApiHandler) RegisterUserHandler(event types.RegisterUser) error {
	if event.Username == "" || event.Password == "" {
		return fmt.Errorf("req has empty parameters")
	}

	userExist, err := api.dbStore.DoesUserExist(event.Username)

	if err != nil {
		return fmt.Errorf("erro occured whne checking user exist %w", err)
	}

	if userExist {
		return fmt.Errorf("user already exist")
	}

	err = api.dbStore.UserInset(event)

	if err != nil {
		return fmt.Errorf("error ocurred while creating a usser %w", err)
	}

	return nil
}
