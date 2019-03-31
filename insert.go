package main

import "./dynamoMethods"

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/ChimeraCoder/anaconda"
	"github.com/aws/aws-lambda-go/lambda"
)

// GoでLambdaを使うときの定型メソッド
func main() {
	log.Println("insert started")
	lambda.Start(insertTweetData)
}

// IncrementID 自動採番用IDの型
type IncrementID int

// Response レスポンス
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// メイン処理メソッド
func insertTweetData(ctx context.Context, params map[string]string) ([]Response, error) {
	// パラメータから検索語を取得
	animalName := params["word"]

	// TwitterAPI用の設定を行う
	api := getTwitterAPI()

	// 各検索ワード毎にツイートを100件(最大値)ずつ取得
	v := url.Values{}
	v.Set("count", "100")

	// 実取得部分
	searchResponse, _ := api.GetSearch(animalName+" filter:native_video -filter:retweets -vine -periscope", v)
	var tweets []anaconda.Tweet
	tweets = searchResponse.Statuses

	fmt.Println("Dynamo書き込み開始")
	// DynamoDBに書き込む
	for _, tweet := range tweets {
		dynamoMethods.WriteDB(animalName, tweet)
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
