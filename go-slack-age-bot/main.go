package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/shomali11/slacker"
)

func loadCredentials() map[string]string {
	content, err := ioutil.ReadFile("./credentials.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	var payload map[string]string
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error when parsing data: ", err)
	}

	return payload
}

func printCommandEvents(analysisticsChannel <-chan *slacker.CommandEvent) {
	for event := range analysisticsChannel {
		fmt.Println("Command Event:")
		fmt.Println(event.Timestamp)
		fmt.Println(event.Command)
		fmt.Println(event.Parameters)
		fmt.Println(event.Event)
	}
}

func main() {
	currentYear := time.Now().Year()
	credentials := loadCredentials()

	// Use Setenv and Getenv to have more potential in extending the application
	// For example: when you use another application to manage the credentials
	os.Setenv("SLACK_BOT_TOKEN", credentials["SLACK_BOT_TOKEN"])
	os.Setenv("SLACK_APP_TOKEN", credentials["SLACK_APP_TOKEN"])
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	go printCommandEvents(bot.CommandEvents())

	bot.Command("my yob is <year>", &slacker.CommandDefinition{
		Description: "yob caculator",
		Examples: []string{"my yob is 2020"},
		Handler: func(botCtx slacker.BotContext, request slacker.Request, responese slacker.ResponseWriter) {
			year := request.Param("year")
			yob, err := strconv.Atoi(year)
			if err != nil {
				fmt.Println("Error at year parameter")
			}

			age := currentYear - yob
			r := fmt.Sprintf("Your age is %d", age)
			responese.Reply(r)
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
