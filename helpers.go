package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

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

type paramParserFunc func(string) (interface{}, error)

func getParams(urlValues url.Values, permitted map[string]paramParserFunc) map[string]interface{} {
	outputMap := make(map[string]interface{})
	for key, f := range permitted {
		if value, ok := urlValues[key]; ok {
			v, err := f(value[0])

			if err == nil {
				outputMap[key] = v
			}
		}
	}

	return outputMap
}
func defaultParamParser(s string) (interface{}, error) {
	return s, nil
}

func booleanParamParser(s string) (interface{}, error) {
	return strconv.ParseBool(s)
}
