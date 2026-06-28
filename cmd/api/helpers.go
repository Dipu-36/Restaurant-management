package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
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

// readJSON method is used fo reading the JSON that the client is trying to send us via POST
// will be used in scnearios like where the client is trying to feedback or their Personal Info
// like their Address, Phone number, Name etc or customizations or cutlery they want
func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {

	// to prevent the APIs from DOS attacks a max read limit of 1MB has been set
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	//initilaize the decoder
	dec := json.NewDecoder(r.Body)
	// calling DisallowUnknowFields methods helps us to prevent if the JSON from the client
	// includes any fields which cannot be mapped to the target destination fo eg trying to
	// write JSON fields which it is not supposed to write, if the client tries something
	// like this the decoder will return an error instead of ignoring the field
	dec.DisallowUnknownFields()

	// decode the request body into the destination
	err := dec.Decode(dst)

	if err != nil {
		var syntaxError *json.SyntaxError                     // struct therefore errors.As()
		var unmarshalTypeError *json.UnmarshalTypeError       // struct therefore again errors.As()
		var invalidUnmarshalError *json.InvalidUnmarshalError // not a struct therefore errors.Is()

		switch {
		// CLient sent malformed JSON like missing brackets in out payloads
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly formed-JSON (at character %d)", syntaxError.Offset)

		// Also malformed JSON but caught differently by the stdlib (Known Golang issue)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly formed JSON")

		// Client sent correct feilds but wrong type like table_number: four instead of 4
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Offset)
			}
			return fmt.Errorf("body contains incorrect JSON type (at charcter %d)", unmarshalTypeError.Offset)

		// Client sent an empty JSON body
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case strings.HasPrefix(err.Error(), "json: unknown field"):
			fieldName := strings.TrimPrefix(err.Error(), "json: unkown field")
			return fmt.Errorf("body contains unkown key %s", fieldName)

		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)

		// We mistakenly pass a non-nil pointer to a Decode() it will return invalidUnmarshlaTypeError
		// since this is a programmer mistake we panic instead of sending this error to the handlers
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		}
	}

	// Decoding 2nd time to catch the clients if they are sending multiple JSON vvalues in one body
	err = dec.Decode(&struct{}{})
	if err != nil {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}
