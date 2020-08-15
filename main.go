package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/cfc-servers/cfc_suggestions/storage"
	"github.com/cfc-servers/cfc_suggestions/storage/sqlite"
	"github.com/cfc-servers/cfc_suggestions/webhooks"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

type suggestionsServer struct {
	storage.SuggestionStore

	suggestionsWebhook webhooks.DiscordWebhook
	suggestionsLoggingWebhook webhooks.DiscordWebhook
	config             *suggestionsConfig
}

func main() {
	port := flag.String("port", "4000", "the port to run the http server on")
	configFile := flag.String("config", "cfc_suggestions_config.json", "configuration file location")
	flag.Parse()

	config := loadConfig(*configFile)

	s := suggestionsServer{
		SuggestionStore:    sqlite.NewStore(config.Database),
		suggestionsWebhook: webhooks.Webhook(config.SuggestionsWebhook),
		suggestionsLoggingWebhook: webhooks.Webhook(config.SuggestionsLoggingWebhook),
		config:             config,
	}

	r := mux.NewRouter()
	r.HandleFunc("/suggestions", s.createSuggestionHandler).Methods("POST")
	r.HandleFunc("/suggestions/{id}/send", s.sendSuggestion).Methods("POST")

	addr := ":" + *port
	http.ListenAndServe(addr, r)
}

func (s *suggestionsServer) createSuggestionHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	var newSuggestionData map[string]string
	json.Unmarshal(body, &newSuggestionData)

	owner, ok := newSuggestionData["owner"]
	if !ok {
		errorJsonResponse(w, http.StatusBadRequest, "Failed to provide  an owner")
		return
	}

	suggestion, err := s.Create(owner)
	if err != nil {
		log.Print(err)
		errorJsonResponse(w, http.StatusInternalServerError, "Database error")
	}
	jsonResponse(w, http.StatusCreated, suggestion)
}

func (s *suggestionsServer) sendSuggestion(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	vars := mux.Vars(r)
	var suggestionCreateData suggestionCreate
	suggestion, _ := s.Get(vars["id"])
	if suggestion == nil || suggestion.Identifier == "" {
		errorJsonResponse(w, http.StatusBadRequest, "Invalid suggestion ID")
		return
	}

	json.Unmarshal(body, &suggestionCreateData)

	embed := webhooks.Embed{
		Title: suggestionCreateData.Title,
		Description: suggestionCreateData.Description,
	}

	err := s.suggestionsWebhook.SendEmbed(embed)
	if err != nil {
		errorJsonResponse(w, http.StatusInternalServerError, "Couldn't send message")
		return
	}
	embed.Fields = append(embed.Fields, &webhooks.EmbedField{
		Name:   "Suggestion Author",
		Value:  fmt.Sprintf("<@!%v>", suggestion.Owner),
	})

	s.suggestionsLoggingWebhook.SendEmbed(embed)

	s.Delete(suggestion.Identifier)
	jsonResponse(w, http.StatusOK, map[string]string{
		"status": "success",
	})
}

func jsonResponse(w http.ResponseWriter, statusCode int, obj interface{}) {
	jsonData, _ := json.Marshal(obj)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(jsonData)
}

func errorJsonResponse(w http.ResponseWriter, statusCode int, err string) {
	obj := map[string]string{"error": err}
	jsonResponse(w, statusCode, obj)
}

type suggestionCreate struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}
