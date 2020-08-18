package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type iEvent struct {
	Id       string `json:"id"`
	Obj_Name string `json:"obj_name"`
	Body     string `json:"body"`
}

func HandleRequest(ctx context.Context, body iEvent) (string, error) {
	fmt.Println(body.Body)
	fmt.Println(body.Id)
	SaveToDynamo(body)
	return fmt.Sprintf(body.Id), nil
}

	func SaveToDynamo(msg iEvent) {
		sess, err := session.NewSession()
		svc := dynamodb.New(sess)

		av, err := dynamodbattribute.MarshalMap(msg)

		if err != nil {
			fmt.Println("Got error marshalling map:")
			fmt.Println(err.Error())
			os.Exit(1)
		}

		input := &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String(os.Getenv("DYNAMO_TABLE_NAME")),
		}

		_, err = svc.PutItem(input)

		if err != nil {
			fmt.Println("Got error calling PutItem:")
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Println("saved to Dynamo " + msg.Id)
	}

func main() {
	lambda.Start(HandleRequest)
}
