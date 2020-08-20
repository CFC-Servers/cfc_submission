package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cfc-servers/cfc_suggestions/suggestions"
	"github.com/cfc-servers/cfc_suggestions/suggestions/sqlite"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlers(t *testing.T) {
	server := suggestionsServer{
		SuggestionStore: sqlite.NewStore("cfc_suggestion_test.db"),
		suggestionsDest: &DummyDest{},
		loggingDest:     &DummyDest{},
		config: &suggestionsConfig{
			Database: "cfc_suggestion_test.db",
		},
	}

	owner1 := "23701337384550"
	owner2 := "247012337384550"

	newSuggestion, _ := server.Create(&suggestions.Suggestion{Owner: owner1})
	newSuggestionForDelete, _ := server.Create(&suggestions.Suggestion{Owner: owner2})
	r := mux.NewRouter()
	r.HandleFunc("/suggestions", server.createSuggestionHandler).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/suggestions/{id}", server.getSuggestionHandler).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/suggestions/{id}", server.deleteSuggestionHandler).Methods(http.MethodDelete)

	testData := []struct {
		handler        http.Handler
		method         string
		endpoint       string
		body           string
		expectedValues map[string]string
		expectedStatus int
		urlVars        map[string]string
	}{
		{
			r,
			"POST",
			"/suggestions",
			`{"owner": "179237013373845504"}`,
			map[string]string{
				"owner": "179237013373845504",
			},
			http.StatusCreated,
			map[string]string{},
		},
		{
			r,
			"POST",
			"/suggestions",
			`{"owner": ""}`,
			map[string]string{},
			http.StatusBadRequest,
			map[string]string{},
		},
		{
			r,
			"POST",
			"/suggestions",
			``,
			map[string]string{},
			http.StatusBadRequest,
			map[string]string{},
		},
		{
			r,
			"GET",
			fmt.Sprintf("/suggestions/%v", newSuggestion.Identifier),
			``,
			map[string]string{},
			http.StatusOK,
			map[string]string{"id": newSuggestion.Identifier},
		},
		{
			r,
			"GET",
			fmt.Sprintf("/suggestions/%v", "13456"),
			``,
			map[string]string{},
			http.StatusNotFound,
			map[string]string{"id": "13456"},
		},
		{
			r,
			"DELETE",
			fmt.Sprintf("/suggestions/%v", "13456"),
			``,
			map[string]string{},
			http.StatusNotFound,
			map[string]string{"id": "13456"},
		},
		{
			r,
			"DELETE",
			fmt.Sprintf("/suggestions/%v", newSuggestionForDelete.Identifier),
			``,
			map[string]string{},
			http.StatusOK,
			map[string]string{"id": newSuggestionForDelete.Identifier},
		},
	}

	for _, data := range testData {
		t.Run(data.endpoint, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			body := data.body
			bodyReader := bytes.NewReader([]byte(body))
			request := httptest.NewRequest(data.method, data.endpoint, bodyReader)
			mux.SetURLVars(request, data.urlVars)
			data.handler.ServeHTTP(recorder, request)

			bodyBytes, _ := ioutil.ReadAll(recorder.Body)
			var returnedValues map[string]string
			json.Unmarshal(bodyBytes, &returnedValues)

			if recorder.Code != data.expectedStatus {
				t.Errorf("Incorrect status %v expected %v", recorder.Code, data.expectedStatus)
			}
			for k, v := range data.expectedValues {
				responseValue, ok := returnedValues[k]
				if !(ok && responseValue == v) {
					t.Errorf("Incorect key %v", k)
				}
			}
		})
	}
}

type DummyDest struct {
}

func (*DummyDest) Send(*suggestions.Suggestion) (string, error)     { return "", nil }
func (*DummyDest) SendEdit(*suggestions.Suggestion) (string, error) { return "", nil }
func (*DummyDest) Delete(string) error                              { return nil }
