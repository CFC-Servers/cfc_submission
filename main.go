package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/cfc-servers/cfc_suggestions/middleware"
	"github.com/cfc-servers/cfc_suggestions/storage"
	"github.com/cfc-servers/cfc_suggestions/storage/sqlite"
	"github.com/cfc-servers/cfc_suggestions/webhooks"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	host := flag.String("host", "127.0.0.1", "the host to run the http server on")
	port := flag.String("port", "4000", "the port to run the http server on")
	configFile := flag.String("config", "cfc_suggestions_config.json", "configuration file location")
	flag.Parse()

	config := loadConfig(*configFile)

	s := suggestionsServer{
		SuggestionStore:           sqlite.NewStore(config.Database),
		suggestionsWebhook:        webhooks.Webhook(config.SuggestionsWebhook),
		suggestionsLoggingWebhook: webhooks.Webhook(config.SuggestionsLoggingWebhook),
		config:                    config,
	}

	r := mux.NewRouter()

	r.Handle(
		"/suggestions",
		middleware.RequireAuth(config.AuthToken, http.HandlerFunc(s.createSuggestionHandler)),
	).Methods(http.MethodPost, http.MethodOptions)

	r.HandleFunc("/suggestions/{id}/send", s.sendSuggestionHandler).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/suggestions/{id}", s.getSuggestionHandler).Methods(http.MethodGet, http.MethodOptions)


	r.Use(
		middleware.SetHeader("Content-Type", "application/json"),
		middleware.LogRequests,

		// CORS
		middleware.SetHeader("Access-Control-Allow-Origin", "https://cfcservers.org"),
		middleware.SetHeader("Access-Control-Allow-Headers", "*"),
		mux.CORSMethodMiddleware(r),
		middleware.IgnoreMethod(http.MethodOptions),
	)

	addr := *host + ":" + *port
	log.Printf("Listening on %v", addr)
	http.ListenAndServe(addr, r)
}

type suggestionsServer struct {
	storage.SuggestionStore

	suggestionsWebhook        webhooks.DiscordWebhook
	suggestionsLoggingWebhook webhooks.DiscordWebhook
	config                    *suggestionsConfig
}

func (s *suggestionsServer) createSuggestionHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	var newSuggestionData map[string]string
	json.Unmarshal(body, &newSuggestionData)

	owner, ok := newSuggestionData["owner"]
	if !ok {
		errorJsonResponse(w, http.StatusBadRequest, "Failed to provide an owner")
		return
	}

	s.DeleteActive(owner)

	suggestion, err := s.Create(owner)
	if err != nil {
		log.Print(err)
		errorJsonResponse(w, http.StatusInternalServerError, "Database error")
		return
	}
	jsonResponse(w, http.StatusCreated, suggestion)
}

func (s *suggestionsServer) getSuggestionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	suggestion, _ := s.GetActive(vars["id"])
	if suggestion == nil || suggestion.Identifier == "" {
		errorJsonResponse(w, http.StatusNotFound, "Suggestion not found")
		return
	}

	jsonResponse(w, http.StatusOK, suggestion)
}

func (s *suggestionsServer) sendSuggestionHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	vars := mux.Vars(r)

	suggestion, _ := s.GetActive(vars["id"])
	if suggestion == nil || suggestion.Identifier == "" {
		errorJsonResponse(w, http.StatusBadRequest, "Invalid suggestion ID")
		return
	}

	var suggestionCreateData suggestionCreate
	json.Unmarshal(body, &suggestionCreateData)

	embed := suggestionCreateData.GetEmbed(suggestion.Owner)
	err := s.suggestionsWebhook.SendEmbed(embed)
	if err != nil {
		errorJsonResponse(w, http.StatusInternalServerError, "Couldn't send message")
		return
	}

	embed.Fields = append(embed.Fields, &webhooks.EmbedField{
		Name:  "Suggestion Author",
		Value: fmt.Sprintf("<@!%v>", suggestion.Owner),
	})
	s.suggestionsLoggingWebhook.SendEmbed(embed)

	s.Update(suggestion.Identifier, false, suggestionCreateData.JsonString())

	jsonResponse(w, http.StatusOK, map[string]string{
		"status": "success",
	})
}

type suggestionCreate struct {
	Realm     string `json:"realm"`
	Link      string `json:"link"`
	Title     string `json:"title"`
	Why       string `json:"why"`
	WhyNot    string `json:"whyNot"`
	Anonymous bool   `json:"anonymous"`
}

func (suggestion suggestionCreate) JsonString() string {
	data, _ := json.Marshal(suggestion)
	return string(data)
}

func (suggestion suggestionCreate) GetEmbed(owner string) webhooks.Embed {
	description := fmt.Sprintf("**%v**\n\n%v", suggestion.Title, suggestion.Link)

    embed := webhooks.Embed{
		Title:       fmt.Sprintf("%v Suggestion", suggestion.Realm),
		Description: description,
		Fields: []*webhooks.EmbedField{
			{
				Name:  "Why",
				Value: suggestion.Why,
			},
			{
				Name:  "Why Not",
				Value: suggestion.WhyNot,
			},
		},
	}

	if !suggestion.Anonymous {
        embed.Fields = append(embed.Fields, &webhooks.EmbedField{
            Name:  "Suggestion Author",
            Value: fmt.Sprintf("<@!%v>", suggestion.Owner),
        })
	}

	return embed
}

func jsonResponse(w http.ResponseWriter, statusCode int, obj interface{}) {
	jsonData, _ := json.Marshal(obj)
	w.WriteHeader(statusCode)
	w.Write(jsonData)
}

func errorJsonResponse(w http.ResponseWriter, statusCode int, err string) {
	obj := map[string]string{"error": err}
	jsonResponse(w, statusCode, obj)
}
