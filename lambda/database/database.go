package database

import (
	"fmt"
	"lambda-func/types"
	"lambda-func/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type UserStore interface {
	DoesUserExist(username string) (bool, error)
	GetUser(username string) (types.LoginUser, error)
	UserInset(user types.User) error
}

type DynamoDBClient struct {
	dbStore *dynamodb.DynamoDB
}

func NewDynamoDBClient() DynamoDBClient {
	dbSession := session.Must(session.NewSession())
	db := dynamodb.New(dbSession)

	return DynamoDBClient{
		dbStore: db,
	}
}

func (u DynamoDBClient) DoesUserExist(username string) (bool, error) {
	result, err := u.dbStore.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(utils.USER_TABLE),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(username),
			},
		},
	})

	if err != nil {
		return true, err
	}

	if result.Item == nil {
		return false, nil
	}

	return true, nil
}

func (u DynamoDBClient) UserInset(user types.User) error {
	item := &dynamodb.PutItemInput{
		TableName: aws.String(utils.USER_TABLE),
		Item: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(user.Username),
			},
			"password": {
				S: aws.String(user.PasswordHash),
			},
		},
	}

	_, err := u.dbStore.PutItem(item)

	if err != nil {
		return err
	}

	return nil
}

func (u DynamoDBClient) GetUser(username string) (types.LoginUser, error) {
	var user types.LoginUser

	result, err := u.dbStore.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(utils.USER_TABLE),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(username),
			},
		},
	})

	if err != nil {
		return user, err
	}

	if result.Item == nil {
		return user, fmt.Errorf("user not found")
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &user)

	if err != nil {
		return user, err
	}

	return user, nil

}
