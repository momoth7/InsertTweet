package anacondaMethods

import (
	"net/url"
	"os"

	"github.com/ChimeraCoder/anaconda"
)

// TwitterAPIから動画付きツイートを取得する
func GetTweets(animalName string) (tweets []anaconda.Tweet) {
	// TwitterAPI用の設定を行う
	api := getTwitterAPI()

	// 各検索ワード毎にツイートを100件(最大値)ずつ取得
	v := url.Values{}
	v.Set("count", "100")

	// 実取得部分
	searchResponse, _ := api.GetSearch(animalName+" filter:native_video -filter:retweets -vine -periscope", v)
	return searchResponse.Statuses
}

// Lambdaの環境変数に基づいてTwitterAPIのセッティング
func getTwitterAPI() *anaconda.TwitterApi {
	anaconda.SetConsumerKey(os.Getenv("CONSUMER_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("CONSUMER_SECRET"))
	return anaconda.NewTwitterApi(os.Getenv("ACCESS_TOKEN"), os.Getenv("ACCESS_TOKEN_SECRET"))
}
