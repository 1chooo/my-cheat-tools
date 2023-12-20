// Author: @1chooo (Hugo ChunHo Lin)
// Created Date: 2023/12/20
// Version: v0.0.1

package main

import (
	"context"
	// "encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/line/line-bot-sdk-go/linebot"
)

func lambdaHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	channelSecret := os.Getenv("CHANNEL_SECRET")
	channelAccessToken := os.Getenv("CHANNEL_ACCESS_TOKEN")

	bot, err := linebot.New(channelSecret, channelAccessToken)
	if err != nil {
		log.Printf("Error creating Line bot: %s", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error creating Line bot",
		}, err
	}

	events, err := bot.ParseRequest(req.Body)
	if err != nil {
		log.Printf("Error parsing request: %s", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Error parsing request",
		}, err
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				handleTextMessage(bot, event.ReplyToken, message.Text)
			default:
				log.Printf("Unknown message: %v", message)
			}
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Message processed successfully",
	}, nil
}

func handleTextMessage(bot *linebot.Client, replyToken, text string) {
	var replyMessages []linebot.SendingMessage

	switch text {
	case "Hello":
		replyMessages = append(replyMessages, linebot.NewTextMessage("World"))
	case "我想看帥哥", "Image":
		replyMessages = append(replyMessages,
			linebot.NewTextMessage("Test get image from s3 public bucket"),
			linebot.NewTextMessage("This is Hugo!"),
			linebot.NewImageMessage("https://2023-amazon-ambassador.s3.amazonaws.com/hugo_grad.png", "https://2023-amazon-ambassador.s3.amazonaws.com/hugo_grad.png"),
		)
	default:
		replyMessages = append(replyMessages, linebot.NewTextMessage(text))
	}

	if _, err := bot.ReplyMessage(replyToken, replyMessages...).Do(); err != nil {
		log.Printf("Error replying message: %s", err)
	}
}

func main() {
	lambda.Start(lambdaHandler)
}
