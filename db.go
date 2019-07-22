package main

import (
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type dynamoDB struct {
	client     *dynamodb.DynamoDB
	tableName  *string
	attribute  *string
	primaryKey *string
}

// NewDynamoDB returns an instance of DynamoDB
func NewDynamoDB(xpub string) *dynamoDB {
	db := new(dynamoDB)
	db.client = dynamodb.New(session.New(), aws.NewConfig().WithRegion("eu-west-1"))
	db.tableName = aws.String("BtcAddressGenerator")
	db.attribute = aws.String("AddressCounter")
	db.primaryKey = aws.String(xpub)
	return db
}

func (db *dynamoDB) ItemExists() bool {
	input := &dynamodb.GetItemInput{
		TableName: db.tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"Xpub": {
				S: db.primaryKey,
			},
		},
	}

	resp, err := db.client.GetItem(input)
	if err != nil {
		return false
	}

	return resp.Item != nil
}

func (db *dynamoDB) AddItem(index string) error {
	input := &dynamodb.PutItemInput{
		TableName: db.tableName,
		Item: map[string]*dynamodb.AttributeValue{
			"Xpub": {
				S: db.primaryKey,
			},
			"AddressCounter": {
				N: aws.String(index),
			},
		},
	}

	_, err := db.client.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}

func (db *dynamoDB) GetCounter() (int, error) {
	input := &dynamodb.GetItemInput{
		TableName: db.tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"Xpub": {
				S: db.primaryKey,
			},
		},
	}

	resp, err := db.client.GetItem(input)
	if err != nil {
		return -1, err
	}

	if resp.Item == nil {
		return -1, nil
	}

	out := map[string]string{}
	err = dynamodbattribute.UnmarshalMap(resp.Item, &out)
	if err != nil {
		return -1, err
	}

	counter, _ := strconv.Atoi(out["AddressCounter"])
	return counter, nil
}

func (db *dynamoDB) IncrementCounter() error {
	input := &dynamodb.UpdateItemInput{
		TableName: db.tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"Xpub": {
				S: db.primaryKey,
			},
		},
		UpdateExpression: aws.String("set AddressCounter = AddressCounter + :val"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":val": {
				N: aws.String("1"),
			},
		},
	}

	_, err := db.client.UpdateItem(input)
	if err != nil {
		return err
	}

	return nil
}
