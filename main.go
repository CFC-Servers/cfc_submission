package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"flag"
	"github.com/bwmarrin/discordgo"
	"github.com/cfc-servers/cfc_suggestions/discord"
	"github.com/cfc-servers/cfc_suggestions/middleware"
	"github.com/cfc-servers/cfc_suggestions/suggestions"
	"github.com/cfc-servers/cfc_suggestions/suggestions/sqlite"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

func main() {
	host := flag.String("host", "127.0.0.1", "the host to run the http server on")
	port := flag.String("port", "4000", "the port to run the http server on")
	configFile := flag.String("config", "cfc_suggestions_config.json", "configuration file location")
	flag.Parse()

	config := loadConfig(*configFile)

	discordgoSession, err := discordgo.New(config.BotToken)
	if err != nil {
		panic(err)
	}
	s := suggestionsServer{
		suggestionsDest: discord.NewDest(config.SuggestionsChannel, false, discordgoSession),
		loggingDest:     discord.NewDest(config.SuggestionsLoggingChannel, true, discordgoSession),
		SuggestionStore: sqlite.NewStore(config.Database),
		config:          config,
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
	log.Infof("Listening on %v", addr)
	http.ListenAndServe(addr, r)
}

type suggestionsServer struct {
	suggestions.SuggestionStore

	suggestionsDest suggestions.Destination
	loggingDest     suggestions.Destination
	config          *suggestionsConfig
}

func (s *suggestionsServer) createSuggestionHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	var newSuggestionData map[string]string
	json.Unmarshal(body, &newSuggestionData)

	owner, _ := newSuggestionData["owner"]
	if owner == "" {
		errorJsonResponse(w, http.StatusBadRequest, "Failed to provide an owner")
		return
	}

	s.Delete(owner, true)

	suggestion, err := s.Create(owner)
	if err != nil {
		log.Error(err)
		errorJsonResponse(w, http.StatusInternalServerError, "Database error")
		return
	}

	jsonResponse(w, http.StatusCreated, suggestion)
}

func (s *suggestionsServer) getSuggestionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	suggestion, err := s.Get(vars["id"])

	if errors.Is(err, sql.ErrNoRows) {
		errorJsonResponse(w, http.StatusNotFound, "Suggestion not found")
		return
	}

	jsonResponse(w, http.StatusOK, suggestion)
}

func (s *suggestionsServer) sendSuggestionHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	vars := mux.Vars(r)

	suggestion, _ := s.Get(vars["id"])
	if suggestion == nil {
		errorJsonResponse(w, http.StatusBadRequest, "Invalid suggestion ID")
		return
	}

	var suggestionContent suggestions.SuggestionContent
	json.Unmarshal(body, &suggestionContent)
	suggestion.Content = &suggestionContent

	s.loggingDest.Send(suggestion)
	if suggestion.Sent {
		_, err := s.suggestionsDest.SendEdit(suggestion)
		if err != nil {
			errorJsonResponse(w, http.StatusInternalServerError, "Couldnt send your suggestion")
		}
		return
	}

	messageId, err := s.suggestionsDest.Send(suggestion)

	suggestion.Sent = true
	suggestion.MessageID = messageId
	if err != nil {
		errorJsonResponse(w, http.StatusInternalServerError, "Couldnt send your suggestion")
		return
	}
	s.Update(suggestion)

	jsonResponse(w, http.StatusOK, map[string]string{
		"status": "success",
	})
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
