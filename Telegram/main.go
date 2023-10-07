// Author: @1chooo (Hugo ChunHo Lin)
// Created Date: 2023/09/27
// Version: v0.0.1

package main

import (
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "os"
)

func main() {
    bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_API_TOKEN"))
    if err != nil {
        panic(err)
    }
    bot.Debug = true
    updateConfig := tgbotapi.NewUpdate(0)
    updateConfig.Timeout = 60
    updates := bot.GetUpdatesChan(updateConfig)
    for update := range updates {
        go handleUpdate(bot, update)
    }
}

func handleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    text := update.Message.Text
    chatID := update.Message.Chat.ID
    replyMsg := tgbotapi.NewMessage(chatID, text)
    if update.Message.IsCommand() {
        switch update.Message.Command() {
			case "start":
				replyMsg.Text = "Hello " + update.Message.From.FirstName
			case "help":
				replyMsg.Text = "What can I help you?"
			default:
				replyMsg.Text = "No such command!!!"
        }
    }
    _, _ = bot.Send(replyMsg)
}
