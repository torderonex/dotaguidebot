package main

import (
	"log"
	"os"
	"strings"

	db "./database"
	"./parser"
	"./screenshoter"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("5858097214:AAHO72oSDS8kezTdz7WDnaNpIOpDrQrSxqE")
	if err != nil {
		panic(err)
	}

	bot.Debug = true

	// Create a new UpdateConfig struct with an offset of 0. Offsets are used
	// to make sure Telegram knows we've handled previous values and we don't
	// need them repeated.
	updateConfig := tgbotapi.NewUpdate(0)

	// Tell Telegram we should wait up to 30 seconds on each request for an
	// update. This way we can get information just as quickly as making many
	// frequent requests without having to send nearly as many.
	updateConfig.Timeout = 30

	// Start polling Telegram for updates.
	updates := bot.GetUpdatesChan(updateConfig)

	// Let's go through each update that we're getting from Telegram.
	for update := range updates {
		// Telegram can send many types of updates depending on what your Bot
		// is up to. We only want to look at messages for now, so we can
		// discard any other updates.
		if update.Message == nil {
			continue
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		switch update.Message.Command() {
		case "hero":
			go func(update tgbotapi.Update) {
				hero := strings.Replace(update.Message.Text, "/hero ", "", 1)
				roles, url, err := parser.Parser(hero)
				if err != nil {
					msg.Text = "This hero does not exist"
					bot.Send(msg)
					return
				}
				go db.Write(update.Message.Chat.ID, url, hero, roles)
				for _, role := range roles {
					msg.Text = role
					if _, err := bot.Send(msg); err != nil {
						log.Panic(err)
					}
				}
				msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(createbuttons(roles))
				msg.Text = "Choose a role"
				bot.Send(msg)
			}(update)

		case "role":
			go func(update tgbotapi.Update) {
				id := update.Message.Chat.ID
				hero, err := db.Get(id)
				if err != nil {
					msg.Text = "You have to choose a hero"
					bot.Send(msg)
					return
				}
				role := strings.Replace(update.Message.Text, "/role ", "", 1)
				roleNum := find(hero.Roles, role)
				if roleNum == -1 {
					msg.Text = "There is no build on this role"
					bot.Send(msg)
					return
				}
				screenshoter.Screenshot(hero.Url, hero.Name+" "+role+".png", roleNum)
				photoBytes, err := os.ReadFile("./pics/" + hero.Name + " " + role + ".png")
				if err != nil {
					log.Fatal(err)
				}
				photoFileBytes := tgbotapi.FileBytes{
					Name:  "picture",
					Bytes: photoBytes,
				}
				_, err = bot.Send(tgbotapi.NewPhoto(id, photoFileBytes))
				if err != nil {
					log.Fatal(err)
				}
			}(update)
		default:
			msg.Text = "I don't know that command"
			bot.Send(msg)
		}
	}
}

func find(arr []string, str string) int {
	for i, x := range arr {
		if str == x {
			return i + 1
		}
	}
	return -1
}

func createbuttons(roles []string) []tgbotapi.KeyboardButton {
	var res []tgbotapi.KeyboardButton
	for _, x := range roles {
		res = append(res, tgbotapi.NewKeyboardButton("/role "+x))
	}
	return res
}
