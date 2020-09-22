package main

import (
	"flag"
	"github.com/bwmarrin/discordgo"
	"github.com/cfc-servers/cfc_suggestions/discord"
	"github.com/cfc-servers/cfc_suggestions/middleware"
	"github.com/cfc-servers/cfc_suggestions/suggestions/sqlite"
	"github.com/getsentry/sentry-go"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func main() {
	host := flag.String("host", "127.0.0.1", "the host to run the http server on")
	port := flag.String("port", "4000", "the port to run the http server on")
	configFile := flag.String("config", "cfc_suggestions_config.json", "configuration file location")
	flag.Parse()

	config := loadConfig(*configFile)

	initSentry(config.SentryDSN)
	defer sentry.Flush(5 * time.Second)

	discordgoSession, err := discordgo.New(config.BotToken)
	if err != nil {
		panic(err)
	}

	sqliteStore := sqlite.NewStore(config.Database)
	sqliteStore.LogQueries = config.LogSql
	s := suggestionsServer{
		suggestionsDest: discord.NewDest(config.SuggestionsChannel, false, discordgoSession),
		loggingDest:     discord.NewDest(config.SuggestionsLoggingChannel, true, discordgoSession),
		SuggestionStore: sqliteStore,
		config:          config,
	}

	r := mux.NewRouter()

	var createSuggestionsHandler http.Handler = http.HandlerFunc(s.createSuggestionHandler)
	var indexSuggestionsHandler http.Handler = http.HandlerFunc(s.indexSuggestionHandler)

	if config.IgnoreAuth {
		log.Warning("RUNNING WITHOUT AUTHENTICATION!!!")
	} else {
		createSuggestionsHandler = middleware.RequireAuth(config.AuthToken, createSuggestionsHandler)
		indexSuggestionsHandler = middleware.RequireAuth(config.AuthToken, indexSuggestionsHandler)
	}

	r.Handle("/suggestions", createSuggestionsHandler).Methods(http.MethodPost, http.MethodOptions)
	r.Handle("/suggestions", indexSuggestionsHandler).Methods(http.MethodGet)

	r.HandleFunc("/suggestions/{id}", s.deleteSuggestionHandler).Methods(http.MethodDelete)
	r.HandleFunc("/suggestions/{id}/send", s.sendSuggestionHandler).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/suggestions/{id}", s.getSuggestionHandler).Methods(http.MethodGet, http.MethodOptions)

	r.Use(
		middleware.Recover,
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
	err = http.ListenAndServe(addr, r)
	if err != nil {
		sentry.CaptureException(err)
		log.Error(err)
	}
}

func initSentry(dsn string) {
	err := sentry.Init(sentry.ClientOptions{
		Dsn: dsn,
	})

	if err != nil {
		log.Fatalf("initSentry: %v", err)
	}
}
