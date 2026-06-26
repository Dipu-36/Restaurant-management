package main

import (
	"encoding/json"
	"net/http"
)

type envelope map[string]interface{}

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {

	// MarshlIndent is a fucntion is to convert Go's data into pretty-printed json
	// 1st param is the data to be converted, 2nd param is the prefix string & the 3rd pram is indentation string
	// prefix string is used to prefix the ouput often used for embedding json inside logs, debugging ouput etc
	// indentation string controls the indentation for each level of nesting \t is for the tabs
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}
