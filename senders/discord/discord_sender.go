package discord

import (
	"github.com/cfc-servers/cfc_suggestions/forms"
	"github.com/go-resty/resty/v2"
)

const editedAtTimeFormat = "Mon Jan _2 15:04:05 2006"

var realms = map[string]string{
	"cfc3":    "Build/Kill",
	"cfcrp":   "DarkRP",
	"cfcmc":   "Minecraft",
	"cfcrvr":  "Raft V Raft",
	"discord": "Discord",
	"other":   "Other",
}

type DiscordSender struct {
	WebhookUrl string
	client *resty.Client
}

func New(webhook string) (*DiscordSender, error) {
	return &DiscordSender{
		WebhookUrl: webhook,
		client:  resty.New(),
	}, nil
}

func (sender *DiscordSender) Send(submission forms.Submission) (string, error) {
	embed := MessageEmbed{
		Title: submission.Fields.Get("title"),
		Description: submission.Fields.Get("description"),
		Image:       &MessageEmbedImage{
			URL:      submission.Fields.Get("image"),
			ProxyURL: "",
			Width:    0,
			Height:   0,
		},
	}

	for k, _ := range submission.Fields {
		if k[0] == '_' {
			continue
		}
		v := submission.Fields.Get(k)
		embed.Fields = append(embed.Fields, &MessageEmbedField{
			Name:   k,
			Value:  v,
			Inline: false,
		})
	}


	var msg Message

	_, err := sender.client.R().
		SetBody(WebhookParams{Embeds: []*MessageEmbed{&embed}}).
		SetResult(msg).
		Post(sender.WebhookUrl)

	if err != nil {
		return "" , err
	}

	return msg.ID, nil
}
