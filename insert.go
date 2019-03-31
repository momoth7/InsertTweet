package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/ChimeraCoder/anaconda"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/guregu/dynamo"
)

type IncrementId int

func main() {
	log.Println("insert started")
	lambda.Start(insertTweetData)
}

// Response レスポンス
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func insertTweetData(ctx context.Context, params map[string]string) ([]Response, error) {
	// パラメータから検索語を取得
	animalName := params["word"]

	// TwitterAPI用の設定を行う
	api := getTwitterAPI()

	// 各検索ワード毎にツイートを100件(最大値)ずつ取得
	v := url.Values{}
	v.Set("count", "100")

	// 実取得部分
	searchResponse, _ := api.GetSearch(params["word"]+" filter:native_video -filter:retweets -vine -periscope", v)
	var tweets []anaconda.Tweet
	tweets = searchResponse.Statuses

	fmt.Println("Dynamo書き込み開始")
	// DynamoDBに書き込む
	for _, tweet := range tweets {
		writeDB(animalName, tweet)
	}
	fmt.Println("Dynamo書き込み終了")

	return []Response{Response{Code: 200, Message: "DB書き込み正常終了"}}, nil
}

// Lambdaの環境変数に基づいてTwitterAPIのセッティング
func getTwitterAPI() *anaconda.TwitterApi {
	anaconda.SetConsumerKey(os.Getenv("CONSUMER_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("CONSUMER_SECRET"))
	return anaconda.NewTwitterApi(os.Getenv("ACCESS_TOKEN"), os.Getenv("ACCESS_TOKEN_SECRET"))
}

// DynamoObject is a generic type used to describe objects in dynamodb
type DynamoObject map[string]*dynamodb.AttributeValue

// writeDB DynamoDBにデータを書き込み
func writeDB(animalName string, tweet anaconda.Tweet) {
	// 次回書き込み時のインクリメントIDを取得
	id := getIncrementID()
	fmt.Println(id)

	// 書き込む
	putData(animalName, tweet, id)

	// インクリメントテーブルを更新
	updateIncrement(id)
}

// IDData JSONデコード用に構造体定義
type IDData struct {
	Name          string      `dynamo:"name"`
	CurrentNumber IncrementId `dynamo:"current_number"`
}

// PutData DynamoDBのレコード構造
type PutData struct {
	ID         IncrementId `dynamo:"id"`
	AnimalName string      `dynamo:"animal_name"`
	TweetData  string      `dynamo:"tweet_data"`
}

// getIncrementID DynamoDBにデータを書き込み
func getIncrementID() IncrementId {
	db := dynamo.New(session.New(), &aws.Config{
		Region: aws.String("us-east-1"),
	})
	table := db.Table("sequences")

	var idData IDData
	table.Get("name", "twimalData").One(&idData)

	return idData.CurrentNumber
}

// putData DynamoDBにデータを書き込み
func putData(animalName string, tweet anaconda.Tweet, id IncrementId) {
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
func updateIncrement(id IncrementId) {
	db := dynamo.New(session.New(), &aws.Config{
		Region: aws.String("us-east-1"),
	})
	table := db.Table("sequences")

	// current_numberを+1する
	putData := IDData{Name: "twimalData", CurrentNumber: id + 1}
	table.Put(putData).Run()
}
