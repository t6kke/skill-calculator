package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/t6kke/skill-calculator/internal/auth"
	"github.com/t6kke/skill-calculator/internal/database"
)

func (api_config *apiConfig) handlerLeagueCreate(w http.ResponseWriter, r *http.Request) {
	type newLeagueParameters struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		UserID      int    `json:"user_id"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}
	user_id, err := auth.ValidateJWT(token, api_config.jwt_secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	new_league_parameters := newLeagueParameters{}
	err = decoder.Decode(&new_league_parameters)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to decode authorization parameters", err)
		return
	}

	if new_league_parameters.Title == ""{
		respondWithError(w, http.StatusBadRequest, "Title for the leauge is required", nil)
		return
	}

	league, err := api_config.db.CreateLeageWithUserRelation(database.CreateLeagueParams{
		Title:    new_league_parameters.Title,
		Description: new_league_parameters.Description,
		UserID: user_id,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create league", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, league)
}


func (api_config *apiConfig) handlerLeagueGet(w http.ResponseWriter, r *http.Request) {
	league_id_string := r.PathValue("leagueID")
	league_id, err := strconv.Atoi(league_id_string)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid leauge ID", err)
		return
	}

	leauge, err := api_config.db.GetLeague(league_id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get leauge", err)
		return
	}

	respondWithJSON(w, http.StatusOK, leauge) //original - old before we use temp urls for videos
}


func (api_config *apiConfig) handlerLeaguesGetAllForUser(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}
	user_id, err := auth.ValidateJWT(token, api_config.jwt_secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	leagues, err := api_config.db.GetLeagues(user_id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retreive leagues", err)
		return
	}

	respondWithJSON(w, http.StatusOK, leagues)
}

func (api_config *apiConfig) handlerLeaguesDeleteOne(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}
	user_id, err := auth.ValidateJWT(token, api_config.jwt_secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	league_id_string := r.PathValue("leagueID")
	league_id, err := strconv.Atoi(league_id_string)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid leauge ID", err)
		return
	}

	league, err := api_config.db.GetLeague(league_id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get leauge", err)
		return
	}

	if league.UserID != user_id {
		respondWithError(w, http.StatusForbidden, "Not allowed to delete this league", err)
		return
	}

	//TODO finish the logic
}
