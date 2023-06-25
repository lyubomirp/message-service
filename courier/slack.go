package courier

import (
	"github.com/slack-go/slack"
	"message-service/types"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

type SlackClient struct {
	AuthToken string
	ChannelId string
	Slack     *slack.Client
}

func InitSlackClient() SlackClient {
	token := os.Getenv("SLACK_TOKEN")
	channelID := os.Getenv("SLACK_CHANNEL_ID")

	if os.Getenv("ENV") != "prod" {
		return SlackClient{
			AuthToken: token,
			ChannelId: channelID,
			Slack:     slack.New(token, slack.OptionDebug(true)),
		}
	}

	return SlackClient{
		AuthToken: token,
		ChannelId: channelID,
		Slack:     slack.New(token),
	}
}

func (client SlackClient) SendMessage(content types.Content) error {
	attachment := slack.Attachment{
		Pretext: string(content.Subject),
		Text:    string(content.Content),
		Color:   "#000000",
		Fields: []slack.AttachmentField{
			{
				Title: "Date",
				Value: time.Now().String(),
			},
		},
	}

	_, timestamp, err := client.Slack.PostMessage(
		client.ChannelId,
		slack.MsgOptionText("New message from bot", false),
		slack.MsgOptionAttachments(attachment),
	)

	if err != nil {
		return err
	}

	log.Infof("Slack Message sent successfully at %s", timestamp)
	return nil
}
