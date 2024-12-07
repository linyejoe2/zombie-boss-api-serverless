package database

import (
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var svc *dynamodb.DynamoDB

func init() {
	sess := session.Must(session.NewSession())
	svc = dynamodb.New(sess, aws.NewConfig().WithRegion("ap-northeast-2"))
}

func PreRegister(email string) error {
	timestamp := time.Now().Unix()

	input := &dynamodb.PutItemInput{
		TableName: aws.String("zombie_boss_pre_register"),
		Item: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
			"timestamp": {
				N: aws.String(fmt.Sprintf("%d", timestamp)),
			},
		},
	}

	_, err := svc.PutItem(input)
	if err != nil {
		return fmt.Errorf("failed to pre-register: %v", err)
	}
	return nil
}

func SetCount(count string) error {
	input := &dynamodb.PutItemInput{
		TableName: aws.String("zombie_boss"),
		Item: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String("zombieBossCount"),
			},
			"Count": {
				N: aws.String(count),
			},
		},
	}

	_, err := svc.PutItem(input)
	if err != nil {
		return fmt.Errorf("failed to set count: %v", err)
	}
	return nil
}

func AddCount() error {
	// 取得目前 count 的數值
	count, err := GetCount()
	if err != nil {
		return err
	}

	// 更新 count
	newCount := count + 1
	return SetCount(fmt.Sprintf("%d", newCount))
}

func GetCount() (int, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String("zombie_boss"),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String("zombieBossCount"),
			},
		},
	}

	result, err := svc.GetItem(input)
	if err != nil {
		return 0, fmt.Errorf("failed to get count: %v", err)
	}

	if result.Item == nil {
		return 0, errors.New("no count found")
	}

	var countStruct struct {
		Count int `json:"Count"`
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &countStruct)
	if err != nil {
		return 0, fmt.Errorf("failed to unmarshal result: %v", err)
	}

	return countStruct.Count, nil
}
