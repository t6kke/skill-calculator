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

func (api_config *apiConfig) handlerGetAllTournamentsInPublicLeague(w http.ResponseWriter, r *http.Request) {
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
		ReportName:         "report_ListTournaments",
		TournamentIDFilter: "",
	}

	exit_code, output_str := bcs_args.BSCExecution()
	if exit_code != 0 {
		error_message := fmt.Sprintf("exit code: %d", exit_code)
		respondWithError(w, http.StatusInternalServerError, "BSC execution failed", errors.New(error_message))
		return
	}

	type replyStruct struct {
		Name        string `json:"name"`
		Version     string `json:"version"`
		Tournaments []struct {
			ID       int    `json:"id"`
			Name     string `json:"name"`
			Date     string `json:"date"`
			Location string `json:"location"`
		} `json:"tournaments"`
	}

	response := replyStruct{}
	err = json.Unmarshal([]byte(output_str), &response)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to compile response to json format", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (api_config *apiConfig) handlerGetPublicTournamentResults(w http.ResponseWriter, r *http.Request) {
	tournament_id_string := r.PathValue("tournamentID")

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
		ReportName:         "report_TournamentResults",
		TournamentIDFilter: tournament_id_string,
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
		Message string `json:"message"`
		Results []struct {
			Position      int    `json:"position"`
			Team          string `json:"team"`
			GamesTotal    int    `json:"games_total"`
			GamesWon      int    `json:"games_won"`
			PointsFor     int    `json:"points_for"`
			PointsAgainst int    `json:"points_against"`
			PointsDiff    int    `json:"points_diff"`
		} `json:"results"`
	}

	response := replyStruct{}
	err = json.Unmarshal([]byte(output_str), &response)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to compile response to json format", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response)
}
