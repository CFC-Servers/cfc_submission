package app

import (
	"github.com/cfc-servers/cfc_suggestions/forms"
	"github.com/cfc-servers/cfc_suggestions/senders/discord"
	"os"
)

var SuggestionsDestination = forms.Destination{
	Name:   "suggestions",
	Sender: discord.New(os.Getenv("SUGGESTIONS_WEBHOOK")),
}
