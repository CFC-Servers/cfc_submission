package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/cfc-servers/cfc_suggestions/discord"
	"github.com/cfc-servers/cfc_suggestions/middleware"
	"github.com/cfc-servers/cfc_suggestions/suggestions/sqlite"
	"github.com/getsentry/sentry-go"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"time"
)

func main() {
	loadConfig()

	host := viper.GetString("host")
	port := viper.GetString("port")

	initSentry(viper.GetString("sentry-dsn"))
	defer sentry.Flush(5 * time.Second)

	discordgoSession, err := discordgo.New(viper.GetString("bot-token"))
	if err != nil {
		panic(err)
	}

	sqliteStore := sqlite.NewStore(viper.GetString("database-file"))
	sqliteStore.LogQueries = viper.GetBool("loq-sql")
	s := suggestionsServer{
		suggestionsDest: discord.NewDest(viper.GetString("suggestions-channel"), false, discordgoSession),
		loggingDest:     discord.NewDest(viper.GetString("suggestions-logging-channel"), true, discordgoSession),
		SuggestionStore: sqliteStore,
	}

	r := mux.NewRouter()

	var createSuggestionsHandler http.Handler = http.HandlerFunc(s.createSuggestionHandler)
	var indexSuggestionsHandler http.Handler = http.HandlerFunc(s.indexSuggestionHandler)

	if viper.GetBool("ignore-auth") {
		log.Warning("RUNNING WITHOUT AUTHENTICATION!!!")
	} else {
		authToken := viper.GetString("auth-token")
		createSuggestionsHandler = middleware.RequireAuth(authToken, createSuggestionsHandler)
		indexSuggestionsHandler = middleware.RequireAuth(authToken, indexSuggestionsHandler)
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

	addr := host + ":" + port
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
