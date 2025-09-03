package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/t6kke/skill-calculator/internal/auth"
	"github.com/t6kke/skill-calculator/internal/database"
)

func (api_config *apiConfig) handlerUsersCreateOne(w http.ResponseWriter, r *http.Request) {
	type auth_parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	parameters := auth_parameters{}
	err := decoder.Decode(&parameters)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to decode authorization parameters", err)
		return
	}

	if parameters.Password == "" || parameters.Email == "" {
		respondWithError(w, http.StatusBadRequest, "Email and password are required", nil)
		return
	}

	hashed_password, err := auth.HashPassword(parameters.Password)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to crate passowrd hash", err)
		return
	}

	user, err := api_config.db.CreateUser(database.CreateUserParams{
		Email:    parameters.Email,
		Password: hashed_password,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	user.Password = "***" //password masking for response, this is not needed to be sent back to the user

	respondWithJSON(w, http.StatusCreated, user)
}
