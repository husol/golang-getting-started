package repositories

import (
  "fmt"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/service/dynamodb"
  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
  "golang-getting-started/models"
)

type UserRepo struct {
  db *dynamodb.DynamoDB
}

func NewUserRepo() *UserRepo  {
  return &UserRepo{db: GetDBConnection()}
}

func (repo *UserRepo) GetByEmail (email string) *models.User {
  var obj models.User

  result, _ := repo.db.GetItem(&dynamodb.GetItemInput{
    TableName: aws.String("users"),
    Key: map[string]*dynamodb.AttributeValue{
      "email": {
        S: aws.String(email),
      },
    },
  })

  err := dynamodbattribute.UnmarshalMap(result.Item, &obj)
  if err != nil {
    panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
  }

  return &obj
}

func (repo *UserRepo) Create (obj models.User) error {
  av, err := dynamodbattribute.MarshalMap(obj)
  if err != nil {
    return err
  }

  input := &dynamodb.PutItemInput{
    Item:      av,
    TableName: aws.String("users"),
  }

  _, err = repo.db.PutItem(input)
  if err != nil {
    return err
  }

  return nil
}