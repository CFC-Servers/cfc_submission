package discord

import (
	"fmt"
	"github.com/cfc-servers/cfc_suggestions/forms"
	"github.com/go-resty/resty/v2"
)

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

func (sender *DiscordSender) Delete(messageid string) error {
	// TODO check status code ensuring it is 200
	resp, err := sender.client.R().
		Delete(sender.WebhookUrl + "/messages/" + messageid)

	if err != nil {
		return err
	}
	fmt.Println("url", sender.WebhookUrl+"/messages/"+messageid)
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
	content := submission.Content

	embed := MessageEmbed{Color: submission.Content.Color}
	embed.Description = content.Description
	embed.Title = content.Title
	embed.Image = &MessageEmbedImage{
		URL: content.Image,
	}

	if !submission.Fields.GetBool("anonymous") { // TODO should fields only be accessed this way in the formatter?
		embed.Author = &MessageEmbedAuthor{
			Name:    submission.OwnerInfo.Name,
			IconURL: submission.OwnerInfo.Avatar,
		}
	}

	for _, field := range content.Fields {
		embed.Fields = append(embed.Fields, &MessageEmbedField{
			Name:  field.Name,
			Value: field.Value,
		})
	}
	embed.Footer = &MessageEmbedFooter{
		Text: "User ID: "+submission.OwnerID,
	}
	return embed
}
