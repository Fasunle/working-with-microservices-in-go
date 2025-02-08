package main

import (
	"authentication/data"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {

	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.ReadJSON(w, r, &requestPayload)

	if err != nil {
		app.ErrorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	user, err := app.Models.User.GetByEmail(requestPayload.Email)

	if err != nil {
		app.ErrorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	valid, err := data.PasswordMatches(user.Password, requestPayload.Password)

	if err != nil || !valid {
		app.ErrorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	err = app.logRequest("auth", fmt.Sprintf("User %s logged in", user.Email))

	if err != nil {
		app.ErrorJSON(w, err)
		return
	}

	payload := JsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}

	app.WriteJSON(w, http.StatusAccepted, payload)

}

func (app *Config) Register(w http.ResponseWriter, r *http.Request) {

	var requestPayload struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Active    bool   `json:"active"`
	}

	err := app.ReadJSON(w, r, &requestPayload)

	if err != nil {
		app.ErrorJSON(w, errors.New("invalid Payload"), http.StatusBadRequest)
		return
	}

	u := data.User{
		Email:     requestPayload.Email,
		FirstName: requestPayload.FirstName,
		LastName:  requestPayload.LastName,
		Active:    requestPayload.Active,
		Password:  requestPayload.Password,
	}

	if requestPayload.Email == "" || requestPayload.Password == "" || requestPayload.FirstName == "" || requestPayload.LastName == "" {
		app.WriteJSON(w, http.StatusBadRequest, JsonResponse{
			Error:   true,
			Message: fmt.Sprintf("Invalid payload!"),
			Data:    nil,
		})
		return
	}

	found, err := app.Models.User.GetByEmail(requestPayload.Email)

	if found != nil {
		app.WriteJSON(w, http.StatusBadRequest, JsonResponse{
			Error:   true,
			Message: fmt.Sprintf("This email (%s) is already in use", u.Email),
			Data:    nil,
		})
		return
	}

	user, err := app.Models.User.Insert(u)

	if err != nil {
		app.ErrorJSON(w, errors.New("error occurred during signup. Try again later"), http.StatusBadRequest)
		return
	}

	fmt.Printf("Email: %s", requestPayload.Email)

	err = app.logRequest("auth", fmt.Sprintf("User %s logged in", u.Email))

	if err != nil {
		app.ErrorJSON(w, err)
		return
	}

	payload := JsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Registered a new user with the email address %s", u.Email),
		Data:    user,
	}

	app.WriteJSON(w, http.StatusAccepted, payload)

}

func (app *Config) logRequest(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	entry.Name = name
	entry.Data = data

	jsonData, _ := json.MarshalIndent(entry, "", "\t")
	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))

	if err != nil {
		return err
	}

	client := &http.Client{}
	_, err = client.Do(request)

	return err
}
