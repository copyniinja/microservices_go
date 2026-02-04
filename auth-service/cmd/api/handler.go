package main

import (
	"auth-service/data"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		fmt.Println(err)
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// validate the user against the database
	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) Register(w http.ResponseWriter, r *http.Request) {
	// 1. Parse request body
	var requestPayload struct {
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Password  string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, errors.New("invalid request body"), http.StatusBadRequest)
		return
	}

	// 2. Basic validation
	if requestPayload.Email == "" || requestPayload.Password == "" || requestPayload.FirstName == "" || requestPayload.LastName == "" {
		app.errorJSON(w, errors.New("all fields are required"), http.StatusBadRequest)
		return
	}

	// 3. Check if email already exists
	_, err = app.Models.User.GetByEmail(requestPayload.Email)
	if err == nil {
		// user already exists
		app.errorJSON(w, errors.New("email already registered"), http.StatusBadRequest)
		return
	} else if !errors.Is(err, sql.ErrNoRows) {
		// some other DB error
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// 4. Create new user object
	newUser := data.User{
		Email:     requestPayload.Email,
		FirstName: requestPayload.FirstName,
		LastName:  requestPayload.LastName,
		Password:  requestPayload.Password,
		Active:    true}

	// 5. Insert user into DB
	id, err := app.Models.User.Insert(newUser)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// 6. Fetch the inserted user to return
	user, err := app.Models.User.GetOne(id)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// 7. Respond with JSON (without password)
	payload := jsonResponse{
		Error:   false,
		Message: "User registered successfully",
		Data:    user,
	}

	app.writeJSON(w, http.StatusCreated, payload)
}
