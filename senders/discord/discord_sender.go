package discord

import (
	"fmt"
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
	client     *resty.Client
}

func New(webhook string) *DiscordSender {
	return &DiscordSender{
		WebhookUrl: webhook,
		client:     resty.New(),
	}
}

func (sender *DiscordSender) Edit(messageid string, submission forms.Submission) error {
	embed := getEmbed(submission)
	var msg Message

	// TODO check status code ensuring it is 200
	resp, err := sender.client.R().
		SetBody(WebhookParams{Embeds: []*MessageEmbed{&embed}}).
		SetResult(&msg).
		Patch(sender.WebhookUrl + "/messages/" + messageid)

	if err != nil {
		return err
	}
	fmt.Println("url", sender.WebhookUrl+"/messages/"+messageid)
	fmt.Println("messageid", msg.ID)
	fmt.Println("body", string(resp.Body()))
	fmt.Println("status", resp.Status())
	return nil
}

func (sender *DiscordSender) Send(submission forms.Submission) (string, error) {
	embed := getEmbed(submission)
	var msg Message

	// TODO check status code ensuring it is 200
	resp, err := sender.client.R().
		SetBody(WebhookParams{Embeds: []*MessageEmbed{&embed}}).
		SetResult(&msg).
		Post(sender.WebhookUrl + "?wait=true")

	if err != nil {
		return "", err
	}
	fmt.Println("messageid", msg.ID)
	fmt.Println("body", string(resp.Body()))
	fmt.Println("status", resp.Status())
	return msg.ID, nil
}

func getEmbed(submission forms.Submission) MessageEmbed {
	embed := MessageEmbed{}

	for k, _ := range submission.Fields {
		v := submission.Fields.Get(k)

		switch k {
		case "description":
			embed.Description = v
		case "image":
			embed.Image = &MessageEmbedImage{
				URL: submission.Fields.Get("image"),
			}
		default:
			embed.Fields = append(embed.Fields, &MessageEmbedField{
				Name:   k,
				Value:  v,
				Inline: false,
			})
		}
	}

	return embed
}
