package main

import (
	"encoding/json"
	"github.com/cfc-servers/cfc_suggestions/suggestions"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type suggestionsServer struct {
	suggestions.SuggestionStore

	suggestionsDest suggestions.Destination
	loggingDest     suggestions.Destination
	config          *suggestionsConfig
}

func (s *suggestionsServer) createSuggestionHandler(w http.ResponseWriter, r *http.Request) {
	var newSuggestionData map[string]string
	unmarshallBody(r, &newSuggestionData)

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
	suggestion, _ := s.Get(vars["id"])
	if suggestion == nil {
		errorJsonResponse(w, http.StatusNotFound, "Suggestion not found")
		return
	}
	jsonResponse(w, http.StatusOK, suggestion)
}

func (s *suggestionsServer) sendSuggestionHandler(w http.ResponseWriter, r *http.Request) {
	var suggestionContent suggestions.SuggestionContent
	unmarshallBody(r, &suggestionContent)

	vars := mux.Vars(r)

	suggestion, _ := s.Get(vars["id"])
	if suggestion == nil {
		errorJsonResponse(w, http.StatusBadRequest, "Invalid suggestion ID")
		return
	}

	suggestion.Content = &suggestionContent

	if suggestion.Sent {
		_, err := s.suggestionsDest.SendEdit(suggestion)
		if err != nil {
			errorJsonResponse(w, http.StatusInternalServerError, "Couldnt send your suggestion")
			return
		}
		s.loggingDest.Send(suggestion)
		jsonResponse(w, http.StatusOK, map[string]string{
			"status": "success",
		})

		return
	}

	messageId, err := s.suggestionsDest.Send(suggestion)

	suggestion.Sent = true
	suggestion.MessageID = messageId
	if err != nil {
		errorJsonResponse(w, http.StatusInternalServerError, "Couldnt send your suggestion")
		return
	}
	s.loggingDest.Send(suggestion)
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

func unmarshallBody(r *http.Request, obj interface{}) {
	data, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(data, obj)
}
