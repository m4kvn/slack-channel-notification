package main

import (
	"github.com/joho/godotenv"
	"github.com/nlopes/slack"
	"flag"
	"os"
	"log"
	"time"
	"strconv"
	"fmt"
)

func main() {
	flags := loadFlags()
	api := slack.New(flags.SlackBotToken)

	if _, err := api.AuthTest(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	} else {
		log.Println("Slack AuthTest is ok.")
	}

	if channel, err := flags.setChannelId(api); err != nil {
		log.Fatal(err)
		os.Exit(1)
	} else {
		log.Printf("Set channel id of %s to %s.\n",
			flags.ChannelName, flags.ChannelId)
		if !channel.IsMember {
			fmt.Printf("This bot is not joining %s.\n", flags.ChannelName)
			fmt.Printf("Invite this bot to %s.\n", flags.ChannelName)
		}
	}

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.ChannelCreatedEvent:
				postMessage(api, ev, flags)
			default:
			}
		}
	}
}

func postMessage(api *slack.Client, ev *slack.ChannelCreatedEvent, flags Flag) {
	go func() {
		time.Sleep(flags.Delay)
		channel, _ := api.GetChannelInfo(ev.Channel.ID)
		attachment := slack.Attachment{
			Color: "good",
			Text: fmt.Sprintf("New channel <#%s> has been created by <@%s>\n",
				channel.ID, channel.Creator),
		}
		if channel.Purpose.Value != "" {
			attachment.Fields = []slack.AttachmentField{{
				Title: "Purpose",
				Value: channel.Purpose.Value,
			}}
		}
		params := slack.PostMessageParameters{
			AsUser: true,
			Attachments: []slack.Attachment{attachment},
		}
		if flags.ChannelId != "" {
			api.PostMessage(flags.ChannelId, "", params)
		}
	}()
}

type Flag struct {
	SlackBotToken string
	ChannelName   string
	ChannelId     string
	Delay         time.Duration
}

func loadFlags() Flag {
	godotenv.Load()
	slackBotToken := flag.String("token",
		os.Getenv("SLACK_BOT_TOKEN"),
		"Set your slack bot token.")
	channelName := flag.String("channel",
		os.Getenv("SLACK_CHANNEL_NAME"),
		"Set your slack notification channel.")
	delayStr := flag.String("delay",
		os.Getenv("DELAY"),
		"Set slack notification delay (default 5 second).")
	flag.Parse()

	if *slackBotToken == "" {
		log.Println("Slack token is require.")
		os.Exit(1)
	}

	if *channelName == "" {
		*channelName = "random"
		log.Println("Set notification channel to random.")
	}

	delay, err := strconv.Atoi(*delayStr)

	if err != nil {
		log.Println("Set delay to default (5 second)")
		delay = 5
	}

	flags := Flag{
		SlackBotToken: *slackBotToken,
		ChannelName:   *channelName,
		Delay:         time.Duration(delay) * time.Second,
	}

	return flags
}

func (f *Flag) setChannelId(api *slack.Client) (slack.Channel, error) {
	channels, err := api.GetChannels(false)
	if err != nil {
		return slack.Channel{}, err
	}
	for _, channel := range channels {
		if channel.Name == f.ChannelName {
			f.ChannelId = channel.ID
			return channel, nil
		}
	}
	return slack.Channel{}, fmt.Errorf("Channel %s is not found.\n", f.ChannelName)
}
