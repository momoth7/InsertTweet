package dynamoMethods

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/ChimeraCoder/anaconda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/guregu/dynamo"
)

// DynamoObject is a generic type used to describe objects in dynamodb
type DynamoObject map[string]*dynamodb.AttributeValue

// writeDB DynamoDBにデータを書き込み
func WriteDB(animalName string, tweet anaconda.Tweet) {
	// 次回書き込み時のインクリメントIDを取得
	id := getIncrementID()
	fmt.Println(id)

	// 書き込む
	putData(animalName, tweet, id)

	// インクリメントテーブルを更新
	updateIncrement(id)
}

// IncrementID 自動採番用IDの型
type IncrementID int

// IDData JSONデコード用に構造体定義
type IDData struct {
	Name          string      `dynamo:"name"`
	CurrentNumber IncrementID `dynamo:"current_number"`
}

// PutData DynamoDBのレコード構造
type PutData struct {
	ID         IncrementID `dynamo:"id"`
	AnimalName string      `dynamo:"animal_name"`
	TweetData  string      `dynamo:"tweet_data"`
}

// getIncrementID DynamoDBにデータを書き込み
func getIncrementID() IncrementID {
	db := dynamo.New(session.New(), &aws.Config{
		Region: aws.String("us-east-1"),
	})
	table := db.Table("sequences")

	var idData IDData
	table.Get("name", "twimalData").One(&idData)

	return idData.CurrentNumber
}

// putData DynamoDBにデータを書き込み
func putData(animalName string, tweet anaconda.Tweet, id IncrementID) {
	db := dynamo.New(session.New(), &aws.Config{
		Region: aws.String("us-east-1"),
	})
	table := db.Table("twimalData")

	byteTweet, _ := json.Marshal(tweet)
	var buf bytes.Buffer
	buf.Write(byteTweet)

	putData := PutData{ID: id, AnimalName: animalName, TweetData: buf.String()}
	table.Put(putData).Run()
}

// updateIncrement インクリメントテーブルを更新
func updateIncrement(id IncrementID) {
	db := dynamo.New(session.New(), &aws.Config{
		Region: aws.String("us-east-1"),
	})
	table := db.Table("sequences")

	// current_numberを+1する
	putData := IDData{Name: "twimalData", CurrentNumber: id + 1}
	table.Put(putData).Run()
}
