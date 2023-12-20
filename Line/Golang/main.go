// Author: @1chooo (Hugo ChunHo Lin)
// Created Date: 2023/09/27
// Version: v0.0.1

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

// For LINE Friends images, please refer to https://developers.line.biz/media/messaging-api/sticker_list.pdf
const (
	BrownImage string = "https://stickershop.line-scdn.net/stickershop/v1/sticker/52002734/iPhone/sticker_key@2x.png"
	ConyImage  string = "https://stickershop.line-scdn.net/stickershop/v1/sticker/52002735/iPhone/sticker_key@2x.png"
	SallyImage string = "https://stickershop.line-scdn.net/stickershop/v1/sticker/52002736/iPhone/sticker_key@2x.png"
	BossImage  string = "https://stickershop.line-scdn.net/stickershop/v1/sticker/51626498/iPhone/sticker_key@2x.png"
)


var bot *linebot.Client


func main() {
	var err error
	bot, err = linebot.New(
		os.Getenv("CHANNEL_SECRET"), 
		os.Getenv("CHANNEL_TOKEN"),
	)
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
				var sendr *linebot.Sender

				// GetMessageQuota: Get how many remain free tier push message quota you still have this month. (maximum 500)
				quota, err := bot.GetMessageQuota().Do()
				if err != nil {
					log.Println("Quota err:", err)
				}
				// message.ID: Msg unique ID
				// message.Text: Msg text
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("msg ID:"+message.ID+":"+"Get:"+message.Text+" , \n OK! remain message:"+strconv.FormatInt(quota.Value, 10))).Do(); err != nil {
					log.Print(err)
				}

				// If user already selected the sender feedback, prepare user nick name and icon here.
				switch {
				case strings.EqualFold(message.Text, "Brown"):
					sendr = linebot.NewSender("Brown", BrownImage)
				case strings.EqualFold(message.Text, "Cony"):
					sendr = linebot.NewSender("Cony", ConyImage)
				case strings.EqualFold(message.Text, "Sally"):
					sendr = linebot.NewSender("Sally", SallyImage)
				case strings.EqualFold(message.Text, "Boss") || strings.Contains(message.Text, "老闆"):
					sendr = linebot.NewSender("Boss", BossImage)
				default:
					// User input other than our provided range, notify the user by quick reply.
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("Please select following LINE Friends to reply you: Brown, Cony, and Sally.")).Do(); err != nil {
						log.Print(err)
					}
				}
				if sendr != nil {
					// Send message with switched sender.
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("Hi, this is "+message.Text+", Nice to meet you.").WithSender(sendr)).Do(); err != nil {
						log.Print(err)
					}
				}

			// Handle only on Sticker message
			case *linebot.StickerMessage:
				var kw string
				for _, k := range message.Keywords {
					kw = kw + "," + k
				}

				outStickerResult := fmt.Sprintf("Received sticker message: %s, pkg: %s kw: %s  text: %s", message.StickerID, message.PackageID, kw, message.Text)
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(outStickerResult)).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}
}
