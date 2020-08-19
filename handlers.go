package main

import (
	"github.com/cfc-servers/cfc_suggestions/suggestions"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
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

	s.DeleteWhere(map[string]interface{}{
		"owner": owner,
		"sent":  true,
	})

	suggestion, err := s.Create(&suggestions.Suggestion{
		Owner: owner,
	})

	if err != nil {
		log.Error(err)
		errorJsonResponse(w, http.StatusInternalServerError, "Database error")
		return
	}

	jsonResponse(w, http.StatusCreated, suggestion)
}

func (s *suggestionsServer) getSuggestionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	foundSuggestions, _ := s.GetWhere(map[string]interface{}{
		"identifier": vars["id"],
	})

	if len(foundSuggestions) == 0 {
		errorJsonResponse(w, http.StatusNotFound, "Suggestion not found")
		return
	}
	jsonResponse(w, http.StatusOK, foundSuggestions[0])
}

func booleanParamParser(s string) (interface{}, error) {
	return strconv.ParseBool(s)
}

func (s *suggestionsServer) indexSuggestionHandler(w http.ResponseWriter, r *http.Request) {
	params := getParams(r.URL.Query(), map[string]paramParserFunc{
		"owner":      defaultParamParser,
		"sent":       booleanParamParser,
		"message_id": defaultParamParser,
	})

	outputSuggestions, _ := s.GetWhere(params)

	jsonResponse(w, http.StatusOK, outputSuggestions)
}

func (s *suggestionsServer) deleteSuggestionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	suggestionsForDelete, _ := s.GetWhere(map[string]interface{}{
		"identifier": vars["id"],
	})
	if len(suggestionsForDelete) == 0 {
		errorJsonResponse(w, http.StatusNotFound, "Suggestion not found")
		return
	}
	suggestion := suggestionsForDelete[0]

	s.suggestionsDest.Delete(suggestion.MessageID)
	err := s.DeleteWhere(map[string]interface{}{
		"identifier": suggestion.Identifier,
	})

	if err != nil {
		log.Errorf("Database error %v", err)
		errorJsonResponse(w, http.StatusInternalServerError, "Database error")
		return
	}

	jsonResponse(w, http.StatusOK, map[string]string{
		"status": "success",
	})
}

func (s *suggestionsServer) sendSuggestionHandler(w http.ResponseWriter, r *http.Request) {
	var suggestionContent suggestions.SuggestionContent
	unmarshallBody(r, &suggestionContent)

	vars := mux.Vars(r)

	foundSuggestions, _ := s.GetWhere(map[string]interface{}{
		"identifier": vars["id"],
	})
	if len(foundSuggestions) == 0 {
		errorJsonResponse(w, http.StatusBadRequest, "Invalid suggestion ID")
		return
	}
	suggestion := foundSuggestions[0]

	suggestion.Content = &suggestionContent

	if suggestion.Sent {
		_, err := s.suggestionsDest.SendEdit(suggestion)
		if err != nil {
			errorJsonResponse(w, http.StatusInternalServerError, "Couldnt send your suggestion")
			return
		}
		s.loggingDest.Send(suggestion)
		s.Update(suggestion)
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