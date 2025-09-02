package main

import (
	"encoding/json"
	"net/http"

	"github.com/t6kke/skill-calculator/internal/database"
)

func (api_config *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type auth_parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	parameters := auth_parameters{}
	err := decoder.Decode(&parameters)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to decode authorization parameters", err)
	}

	if parameters.Password == "" || parameters.Email == "" {
		respondWithError(w, http.StatusBadRequest, "Email and password are required", nil)
		return
	}

	user, err := api_config.db.CreateUser(database.CreateUserParams{
		Email:    parameters.Email,
		Password: parameters.Password,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, user) //TODO do I need to respond with full user struct where password hash is in the content
}
