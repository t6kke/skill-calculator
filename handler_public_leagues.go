package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/t6kke/skill-calculator/internal/bsc"
)

func (api_config *apiConfig) handlerGetAllPublicLeagues(w http.ResponseWriter, r *http.Request) {
	leagues, err := api_config.db.GetPublicLeagues()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retreive leagues", err)
		return
	}
	respondWithJSON(w, http.StatusOK, leagues)
}

func (api_config *apiConfig) handlerGetPublicLeagueSandings(w http.ResponseWriter, r *http.Request) {
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
	if !league.IsPublic {
		respondWithError(w, http.StatusForbidden, "Not publicly available league", err)
		return
	}

	bcs_args := bsc.ExecutionArguments{
		Command:            "report",
		DBName:             filepath.Join(api_config.db_dir, league.DatabaseName),
		ReportName:         "report_EloStandings",
		TournamentIDFilter: "",
	}
	exit_code, output_str := bcs_args.BSCExecution()
	if exit_code != 0 {
		error_message := fmt.Sprintf("exit code: %d", exit_code)
		respondWithError(w, http.StatusInternalServerError, "BSC execution failed", errors.New(error_message))
		return
	}

	type replyStruct struct {
		Name     string `json:"name"`
		Version  string `json:"version"`
		Message  string `json:"message"`
		Category []struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			Ranking     []struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
				Elo  int    `json:"elo"`
			} `json:"ranking"`
		} `json:"category"`
	}

	response := replyStruct{}
	err = json.Unmarshal([]byte(output_str), &response)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to compile response to json format", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response)
}
