package main

import (
	"flag"
	"github.com/bwmarrin/discordgo"
	"github.com/cfc-servers/cfc_suggestions/discord"
	"github.com/cfc-servers/cfc_suggestions/middleware"
	"github.com/cfc-servers/cfc_suggestions/suggestions/sqlite"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
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

	r.Handle(
		"/suggestions",
		http.HandlerFunc(s.indexSuggestionHandler),
	).Methods(http.MethodGet)

	r.HandleFunc("/suggestions/{id}", s.deleteSuggestionHandler).Methods(http.MethodDelete)
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
