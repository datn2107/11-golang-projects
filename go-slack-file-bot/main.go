package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/slack-go/slack"
)

func setCredentials() {
	content, err := ioutil.ReadFile("./credentials.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	var payload map[string]string
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error when parsing data: ", err)
	}

	for k, v := range payload {
		os.Setenv(k, v)
	}
}

func main() {
	setCredentials()
	api := slack.New(os.Getenv("SLACK_BOT_TOKEN"))
	channelArr := []string{os.Getenv("CHANNEL_ID")}
	fileArr := []string{"valid.csv"}

	for i := 0; i < len(fileArr); i++ {
		param := slack.FileUploadParameters{
			Channels: channelArr,
			File:     fileArr[i],
		}

		file, err := api.UploadFile(param)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Name: %s, URL: %s\n", file.Name, file.URLPrivate)
	}
}
