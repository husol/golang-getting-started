package repositories

import (
  "fmt"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/service/dynamodb"
  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
  "github.com/aws/aws-sdk-go/service/dynamodb/expression"
  "golang-getting-started/models"
  "os"
)

type PostRepo struct {
  db *dynamodb.DynamoDB
}

func NewPostRepo() *PostRepo  {
  return &PostRepo{db: GetDBConnection()}
}

func (repo *PostRepo) Find(columns []string, conditions []models.Condition,  paging models.Paging, order string) ([]models.Post, int) {
  expr := expression.NewBuilder()

  if len(columns) > 0 {
    var projection expression.ProjectionBuilder
    for _, column := range columns {
      projection = projection.AddNames(expression.Name(column))
    }
    expr = expr.WithProjection(projection)
  }

  if len(conditions) > 0 {
    var filter expression.ConditionBuilder
    for index, condition := range conditions {
      var cond expression.ConditionBuilder
      switch condition.Operator {
      case ">":
        cond = expression.Name(condition.Field).GreaterThan(expression.Value(condition.Values[0]))
      case ">=":
        cond = expression.Name(condition.Field).GreaterThanEqual(expression.Value(condition.Values[0]))
      case "<":
        cond = expression.Name(condition.Field).LessThan(expression.Value(condition.Values[0]))
      case "<=":
        cond = expression.Name(condition.Field).LessThanEqual(expression.Value(condition.Values[0]))
      case "between":
        cond = expression.Name(condition.Field).Between(expression.Value(condition.Values[0]), expression.Value(condition.Values[1]))
      default:
        cond = expression.Name(condition.Field).Equal(expression.Value(condition.Values[0]))
      }

      if index > 0 {
        filter = filter.And(cond)
      } else {
        filter = cond
      }
    }
    expr = expr.WithFilter(filter)
  }

  exprBuilder, err := expr.Build()

  if err != nil {
    fmt.Println("Got error building expression:")
    fmt.Println(err.Error())
    os.Exit(1)
  }

  // Build the query input parameters
  params := &dynamodb.ScanInput{
    ExpressionAttributeNames:  exprBuilder.Names(),
    ExpressionAttributeValues: exprBuilder.Values(),
    FilterExpression:          exprBuilder.Filter(),
    ProjectionExpression:      exprBuilder.Projection(),
    TableName:                 aws.String("posts"),
  }

  if paging.Size > 0 {
    params.SetLimit(int64(paging.Size))
  }

  // Make the DynamoDB Query API call
  numPage := 0
  var result *dynamodb.ScanOutput
  err = repo.db.ScanPages(params,
    func(page *dynamodb.ScanOutput, lastPage bool) bool {
      if numPage < paging.Index {
        result = page
      }
      numPage++
      params.SetExclusiveStartKey(page.LastEvaluatedKey)

      return !lastPage
    })

  if err != nil {
    fmt.Println("Query API call failed:")
    fmt.Println(err.Error())
    os.Exit(1)
  }

  var results []models.Post
  for _, i := range result.Items {
    item := models.Post{}

    err = dynamodbattribute.UnmarshalMap(i, &item)

    if err != nil {
      fmt.Println("Got error unmarshalling:")
      fmt.Println(err.Error())
      os.Exit(1)
    }

    results = append(results, item)
  }

  return results, numPage
}

func (repo *PostRepo) GetById (id string) *models.Post {
  var obj models.Post

  result, _ := repo.db.GetItem(&dynamodb.GetItemInput{
    TableName: aws.String("posts"),
    Key: map[string]*dynamodb.AttributeValue{
      "id": {
        S: aws.String(id),
      },
    },
  })

  err := dynamodbattribute.UnmarshalMap(result.Item, &obj)
  if err != nil {
    panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
  }

  return &obj
}

func (repo *PostRepo) Create (obj models.Post) error {
  av, err := dynamodbattribute.MarshalMap(obj)
  if err != nil {
    return err
  }

  input := &dynamodb.PutItemInput{
    Item:      av,
    TableName: aws.String("posts"),
  }

  _, err = repo.db.PutItem(input)
  if err != nil {
    return err
  }

  return nil
}

func (repo *PostRepo) Update(obj *models.Post) error {
  input := &dynamodb.UpdateItemInput{
    ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
      ":title": {
        S: aws.String(obj.Title),
      },
      ":content": {
        S: aws.String(obj.Content),
      },
    },
    TableName: aws.String("posts"),
    Key: map[string]*dynamodb.AttributeValue{
      "id": {
        S: aws.String(obj.ID),
      },
    },
    ReturnValues:     aws.String("UPDATED_NEW"),
    UpdateExpression: aws.String("set title = :title, content = :content"),
  }

  _, err := repo.db.UpdateItem(input)
  if err != nil {
    return err
  }

  return nil
}

func (repo *PostRepo) Delete(obj *models.Post) error {
  input := &dynamodb.DeleteItemInput{
    Key: map[string]*dynamodb.AttributeValue{
      "id": {
        S: aws.String(obj.ID),
      },
    },
    TableName: aws.String("posts"),
  }

  _, err := repo.db.DeleteItem(input)
  if err != nil {
    return err
  }

  return nil
}