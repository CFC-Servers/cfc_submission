package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/cfc-servers/cfc_suggestions/app"
	"github.com/cfc-servers/cfc_suggestions/dynamodb"
	"github.com/cfc-servers/cfc_suggestions/forms"
	"github.com/cfc-servers/cfc_suggestions/util"
	"github.com/plally/goslash/goslash"
	"github.com/plally/goslash/listeners/lambda"
)

var suggestionsBaseUrl = os.Getenv("SUGGESTIONS_BASE_URL")
var publicKey = os.Getenv("DISCORD_PUBLIC_KEY")
var botToken = os.Getenv("DISCORD_BOT_TOKEN")
var clientId = os.Getenv("DISCORD_CLIENT_ID")
var guildID = os.Getenv("GUILD_ID")

func main() {
	goslashApp, err := goslash.NewApp(clientId, "Bot "+botToken)
	if err != nil {
		panic(err)
	}

	suggestCommand := goslash.NewCommand("suggest", "Create a new suggestion").
		SetHandler(suggestCommandHandler)
	mySubmissionsCommand := goslash.NewCommand("mysubmissions", "Get a list of your pending or submitted submissions").
		SetHandler(mysubmissionsCommandHandler)

	goslashApp.AddCommand(suggestCommand)
	goslashApp.AddCommand(mySubmissionsCommand)

	if _, exists := os.LookupEnv("REGISTER_DISCORD_COMMANDS"); exists {
		goslashApp.RegisterAllGuild(guildID)
		return
	}
	fmt.Println("starting listener")
	listener := lambda.NewListener(publicKey)
	goslashApp.SetListener(listener)
	listener.Start()

}

func mysubmissionsCommandHandler(ctx *goslash.InteractionUpdate) *goslash.InteractionResponse {
	submissions, err := dynamodb.GetOwnerSubmissions(util.GetTable(), ctx.GetAuthor().ID)
	if err != nil {
		return goslash.Response("There was a problem fetching your submissions").Ephemeral()
	}

	var builder strings.Builder
	builder.WriteString("**Active Submissions**\n\n")
	for _, submission := range submissions {
		if builder.Len() > 1600 {
			break
		}
		builder.WriteString(submission.FormName)
		builder.WriteString("- ")

		if submission.Content.Description == "" {
			builder.WriteString("PENDING")
		}
		title := strings.Split(submission.Content.Description, "\n")[0]
		builder.WriteString(title)
		builder.WriteString(" - ")
		builder.WriteString(suggestionsBaseUrl + submission.UUID)
		builder.WriteByte('\n')
	}
	return goslash.Response(builder.String()).Ephemeral()
}

func suggestCommandHandler(ctx *goslash.InteractionUpdate) *goslash.InteractionResponse {
	form, err := app.GetForm("suggestion")
	if err != nil {
		return goslash.Response("sorry there was an error creating your suggestion").Ephemeral()
	}

	avatar := fmt.Sprintf("https://cdn.discordapp.com/avatars/%v/%v.png?size=1024", ctx.GetAuthor().ID, ctx.GetAuthor().Avatar)

	existingSubmissions, _ := dynamodb.GetOwnerSubmissions(util.GetTable(), ctx.GetAuthor().ID)
	for _, submission := range existingSubmissions {
		if submission.FormName != "suggestion" {
			continue
		}
		if len(submission.MessageIDS) == 0 {
			dynamodb.DeleteSubmission(util.GetTable(), submission.UUID)
		}
	}

	submission := forms.NewSubmission(form, forms.OwnerInfo{
		ID:     ctx.GetAuthor().ID,
		Name:   ctx.GetAuthor().Username + "#" + ctx.GetAuthor().Discriminator,
		Avatar: avatar,
		URL:    "",
	})

	_ = dynamodb.PutSubmission(util.GetTable(), submission)
	url := suggestionsBaseUrl + submission.UUID
	return goslash.Response("Click to make a suggestion. **Do not share this URL with anyone!**\n" + url).Ephemeral()
}
