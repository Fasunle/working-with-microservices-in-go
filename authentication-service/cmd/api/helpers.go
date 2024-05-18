package main

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type JsonResponse struct {
	Error   bool   `json: "error"`
	Message string `json: "message"`
	Data    any    `json: "data,omitempty"`
}

// random string source
const randomStringSource = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// reads data as a JSON object and returns an error if the JSON is invalid
func (app *Config) ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {

	const maxBytes = 1024 * 1024 * 1 // 1MB

	r.Body = http.MaxBytesReader(w, r.Body, maxBytes)

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(data)

	if err != nil {
		return err
	}

	err = decoder.Decode(&struct{}{})

	if err != io.EOF {
		return errors.New("body must have only a single JSON value")
	}

	return nil
}

// writes data as a JSON object
func (app *Config) WriteJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)

	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(out)

	if err != nil {
		return err
	}

	return nil

}

// Send error message back to the client
func (app *Config) ErrorJSON(w http.ResponseWriter, err error, status ...int) error {

	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload JsonResponse

	payload.Error = true
	payload.Message = err.Error()

	return app.WriteJSON(w, statusCode, payload)
}

// RandomString returns a random string of letters of length n, using characters specified in randomStringSource.
func (t *Config) RandomString(n int) string {
	s, r := make([]rune, n), []rune(randomStringSource)
	for i := range s {
		p, _ := rand.Prime(rand.Reader, len(r))
		x, y := p.Uint64(), uint64(len(r))
		s[i] = r[x%y]
	}
	return string(s)
}
