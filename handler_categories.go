package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/t6kke/skill-calculator/internal/auth"
	"github.com/t6kke/skill-calculator/internal/bsc"
)

func (api_config *apiConfig) handlerGetCategories(w http.ResponseWriter, r *http.Request) {
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
		respondWithError(w, http.StatusForbidden, "Access rights missing to manage this league", err)
		return
	}

	bcs_args := bsc.ExecutionArguments{
		Command:     "category",
		DBName:      filepath.Join(api_config.db_dir, league.DatabaseName),
		ListContent: true,
	}

	exit_code, output_str := bcs_args.BSCExecution()
	if exit_code != 0 {
		error_message := fmt.Sprintf("exit code: %d", exit_code)
		respondWithError(w, http.StatusInternalServerError, "BSC execution failed", errors.New(error_message))
		return
	}

	type replyStruct struct {
		Name       string `json:"name"`
		Version    string `json:"version"`
		Message    string `json:"message"`
		Categories []struct {
			ID          int    `json:"ID"`
			Name        string `json:"name"`
			Description string `json:"description"`
		} `json:"categories"`
	}

	response := replyStruct{}
	err = json.Unmarshal([]byte(output_str), &response)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to compile response to json format", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (api_config *apiConfig) handlerAddCategory(w http.ResponseWriter, r *http.Request) {
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
		respondWithError(w, http.StatusForbidden, "Access rights missing to manage this league", err)
		return
	}

	type parameters struct {
		CategoryName string `json:"categoryName"`
		CategoryDesc string `json:"categoryDesc"`
	}
	payload := r.FormValue("data")
	params := parameters{}
	err = json.Unmarshal([]byte(payload), &params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to decode input parameters", err)
		return
	}
	//TODO validate that data has been given and not empty values, not an issue from web since fields are required but if some other tool uses the rest endpoint there is no other checks for empty values

	bcs_args := bsc.ExecutionArguments{
		Command:      "category",
		DBName:       filepath.Join(api_config.db_dir, league.DatabaseName),
		CategoryName: params.CategoryName,
		CategoryDesc: params.CategoryDesc,
	}

	exit_code, output_str := bcs_args.BSCExecution()
	if exit_code != 0 {
		error_message := fmt.Sprintf("exit code: %d", exit_code)
		respondWithError(w, http.StatusInternalServerError, "BSC execution failed", errors.New(error_message))
		return
	}

	type replyStruct struct {
		Name    string `json:"name"`
		Version string `json:"version"`
		Status  string `json:"status"`
		Message string `json:"message"`
	}

	response := replyStruct{}
	err = json.Unmarshal([]byte(output_str), &response)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to compile response to json format", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response)
}
