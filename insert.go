package main

import "./dynamoMethods"
import "./anacondaMethods"

import (
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
)

// GoでLambdaを使うときの定型メソッド
func main() {
	log.Println("insert started")
	lambda.Start(lambdaMain)
}

// Response レスポンス
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Request リクエストパラメータの型
type Request struct {
	AnimalName string `json:"animal_name"`
}

// メイン処理メソッド
func lambdaMain(req Request) ([]Response, error) {

	// バリデーション
	if req.AnimalName == "" {
		return []Response{Response{Code: 400, Message: "リクエストパラメータが不正です"}}, nil
	}

	// ツイート情報取得
	tweets := anacondaMethods.GetTweets(req.AnimalName)

	// Dynamoにinsert
	fmt.Println("Dynamo書き込み開始")
	for _, tweet := range tweets {
		dynamoMethods.WriteDB(req.AnimalName, tweet)
	}
	fmt.Println("Dynamo書き込み終了")

	return []Response{Response{Code: 200, Message: "DB書き込み正常終了"}}, nil
}
