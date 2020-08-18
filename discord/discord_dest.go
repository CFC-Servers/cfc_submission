package discord

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/cfc-servers/cfc_suggestions/suggestions"
	log "github.com/sirupsen/logrus"
	"time"
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

type DiscordDestination struct {
	session        *discordgo.Session
	channelId      string
	loggingChannel bool
}

func NewDest(channelId string, isLoggingChannel bool, session *discordgo.Session) *DiscordDestination {
	return &DiscordDestination{
		session:        session,
		channelId:      channelId,
		loggingChannel: isLoggingChannel,
	}
}

func (dest *DiscordDestination) Send(suggestion *suggestions.Suggestion) (string, error) {
	embed := dest.getEmbed(suggestion)
	message, err := dest.session.ChannelMessageSendEmbed(dest.channelId, embed)

	if err != nil {
		err = fmt.Errorf("Error sending destination: %w", err)
		log.Error(err)
		return "", err
	}

	return message.ID, err
}

func (dest *DiscordDestination) Delete(messageId string) error {
	return dest.session.ChannelMessageDelete(dest.channelId, messageId)
}

func (dest *DiscordDestination) SendEdit(suggestion *suggestions.Suggestion) (string, error) {

	if suggestion.MessageID == "" {
		return "", errors.New("Invalid message id")
	}
	embed := dest.getEmbed(suggestion)
	embed.Footer = &discordgo.MessageEmbedFooter{
		Text: fmt.Sprintf("Edited at %v", time.Now().Format(editedAtTimeFormat)),
	}

	message, err := dest.session.ChannelMessageEditEmbed(dest.channelId, suggestion.MessageID, embed)

	if err != nil {
		err = fmt.Errorf("Error sending destination: %w", err)
		log.Error(err)
		return "", err
	}

	return message.ID, err
}

func (dest *DiscordDestination) getEmbed(suggestion *suggestions.Suggestion) *discordgo.MessageEmbed {
	content := suggestion.Content
	humanFriendlyRealm, ok := realms[content.Realm]
	if !ok {
		humanFriendlyRealm = "Other"
	}

	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("[%v] Suggestion", humanFriendlyRealm),
		Description: fmt.Sprintf("**__%v__**\n\n%v", content.Title, content.Link),
		Color:       0x34EB77,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Why?",
				Value:  content.Why,
				Inline: false,
			},
			{
				Name:   "Why Not?",
				Value:  content.WhyNot,
				Inline: false,
			},
		},
	}
	if dest.loggingChannel || !content.Anonymous {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "Author",
			Value:  fmt.Sprintf("<@!%v>", suggestion.Owner),
			Inline: false,
		})
	}
	if dest.loggingChannel {
		if len(suggestion.MessageID) > 0 {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:  "message id",
				Value: suggestion.MessageID,
			})
		}
	}

	return embed
}
