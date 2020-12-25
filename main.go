// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/zhshch2002/goribot"
)

var bot *linebot.Client

func main() {
	var err error
	bot, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	log.Println("Bot:", bot, " err:", err)
	http.HandleFunc("/callback", callbackHandler)
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	events, err := bot.ParseRequest(r)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if err != nil {
					log.Println("Quota err:", err)
				}
				responseHandler(string(message.Text), event.ReplyToken, err)
			}
		}
	}
}

func responseHandler(sign string, replyToken string, err error) {
	var url string

	s := goribot.NewSpider()
	switch sign {
	case "水瓶":
		url = "https://astro.click108.com.tw/daily_10.php?iAstro=10"
	case "雙魚":
		url = "https://astro.click108.com.tw/daily_11.php?iAstro=11"
	case "牡羊":
		url = "https://astro.click108.com.tw/daily_0.php?iAstro=0"
	case "金牛":
		url = "https://astro.click108.com.tw/daily_1.php?iAstro=1"
	case "雙子":
		url = "https://astro.click108.com.tw/daily_2.php?iAstro=2"
	case "巨蟹":
		url = "https://astro.click108.com.tw/daily_3.php?iAstro=3"
	case "獅子":
		url = "https://astro.click108.com.tw/daily_4.php?iAstro=4"
	case "處女":
		url = "https://astro.click108.com.tw/daily_5.php?iAstro=5"
	case "天秤":
		url = "https://astro.click108.com.tw/daily_6.php?iAstro=6"
	case "天蠍":
		url = "https://astro.click108.com.tw/daily_7.php?iAstro=7"
	case "射手":
		url = "https://astro.click108.com.tw/daily_8.php?iAstro=8"
	case "魔羯":
		url = "https://astro.click108.com.tw/daily_9.php?iAstro=9"
	default:
		url = ""
	}

	s.AddTask(
		goribot.GetReq(url),
		func(ctx *goribot.Context) {
			src := ctx.Resp.Text

			//將 HTML 標籤全轉換成小寫
			re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
			src = re.ReplaceAllStringFunc(src, strings.ToLower)

			//去除 STYLE
			re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
			src = re.ReplaceAllString(src, "")

			//去除 SCRIPT
			re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
			src = re.ReplaceAllString(src, "")

			//去除所有尖括號內的 HTML 程式碼，並換成換行符
			re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
			src = re.ReplaceAllString(src, "\n")

			//去除連續的換行符
			re, _ = regexp.Compile("\\s")
			src = re.ReplaceAllString(src, "")

			start := strings.Index(src, "今日"+sign+"座解析")
			end := strings.Index(src, "把屬於你的好運")

			if _, err = bot.ReplyMessage(replyToken, linebot.NewTextMessage(src[start:end])).Do(); err != nil {
				log.Print(err)
			}
		})

	s.Run()
}
