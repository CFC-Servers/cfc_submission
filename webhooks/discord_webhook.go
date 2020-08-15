package webhooks

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type DiscordWebhook struct {
	url    string
	client http.Client
}

func Webhook(url string) DiscordWebhook {
	return DiscordWebhook{url: url, client: http.Client{}}
}

func (webhook *DiscordWebhook) Send(msg Message) error {
	body, _ := json.Marshal(msg)
	bodyReader := bytes.NewReader(body)
	req, _ := http.NewRequest("POST", webhook.url, bodyReader)
	req.Header.Set("Content-Type", "application/json")
	_, err := webhook.client.Do(req)
	return err
}

func (webhook *DiscordWebhook) SendEmbed(embed Embed) error {
	return webhook.Send(Message{
		Content: "",
		Embeds: []Embed{
			embed,
		},
	})
}
