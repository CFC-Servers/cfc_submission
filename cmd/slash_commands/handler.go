package main

import (
	"github.com/cfc-servers/cfc_suggestions/app"
	"github.com/cfc-servers/cfc_suggestions/forms"
	"github.com/plally/goslash/goslash"
	"github.com/plally/goslash/listeners/lambda"
	"os"
)

var baseUrl = os.Getenv("SUGGESTIONS_BASE_URL")
var publicKey = os.Getenv("DISCORD_PUBLIC_KEY")

func main() {
	listener := lambda.NewListener(publicKey)
	goslashApp, err := goslash.NewApp("", "")
	if err != nil {
		panic(err)
	}
	goslashApp.SetListener(listener)

	registerCommands(goslashApp)
	listener.Start()
}

func registerCommands(goslashApp *goslash.Application) {
	suggestCommand := goslash.NewCommand("suggest", "Create a new suggestion")
	suggestCommand.SetHandler("suggest", suggestCommandHandler)
	goslashApp.AddCommand(suggestCommand)
	_, _ = goslashApp.Register("794721717209530368", suggestCommand)

}

func suggestCommandHandler(ctx *goslash.InteractionContext) *goslash.InteractionResponse {
	form, err := app.GetForm("suggestion")
	if err != nil {
		return goslash.Response("sorry there was an error creating your suggestion").OnlyAuthor()
	}
	submission := forms.NewSubmission(form, forms.OwnerInfo{
		ID:     ctx.Member.User.ID,
		Name:   ctx.Member.User.Username+"#"+ctx.Member.User.Discriminator,
		Avatar: ctx.Member.User.Avatar,
		URL:     "",
	})

	url := baseUrl + submission.UUID
	return goslash.Response(url).OnlyAuthor()
}
