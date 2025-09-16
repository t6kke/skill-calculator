package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/t6kke/skill-calculator/internal/auth"
	"github.com/t6kke/skill-calculator/internal/database"
)

func (api_config *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type auth_parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		database.User
		Token string `json:"token"`
	}

	decoder := json.NewDecoder(r.Body)
	parameters := auth_parameters{}
	err := decoder.Decode(&parameters)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to decode authorization parameters", err)
		return
	}

	user, err := api_config.db.GetUserByEmail(parameters.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	err = auth.CheckPasswordHash(user.Password, parameters.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	access_token, err := auth.MakeJWT(user.ID, api_config.jwt_secret, time.Hour*1)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed creating JWT", err)
		return
	}

	user.Password = "***" //password masking for response, this is not needed to be sent back to the user

	respondWithJSON(w, http.StatusOK, response{
		User:  user,
		Token: access_token,
	})
}
