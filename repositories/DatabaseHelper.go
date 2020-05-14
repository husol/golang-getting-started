package repositories

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var dbcon *dynamodb.DynamoDB

func GetDBConnection() *dynamodb.DynamoDB {
	return dbcon
}

func InitDb(config *aws.Config) {
	// Connect dynamodb
	dbcon = dynamodb.New(session.Must(session.NewSession(config)))
}
