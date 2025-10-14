package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/t6kke/skill-calculator/internal/auth"
	"github.com/t6kke/skill-calculator/internal/bsc"
)

//TODO functions
// 1. current standings - in handler_leagues.go - done
// 2. Tournaments in leauge - done
// 3. Tournament reults - done
// 4. Categories in league - in handler_categories.go -
// 5. Add category - in handler_categories.go -

func (api_config *apiConfig) handlerUploadTournament(w http.ResponseWriter, r *http.Request) {
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

	const maxMemory = 10 << 22 // 40 MB using bit shifting
	r.ParseMultipartForm(maxMemory)

	file, header, err := r.FormFile("excel")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to parse form file", err)
		return
	}
	defer file.Close()

	media_type := header.Header.Get("Content-Type")

	m_type, _, err := mime.ParseMediaType(media_type)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to parse the filetype of provided file", err)
		return
	}
	//log.Printf("Media Type: %s", m_type)
	if m_type != "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" {
		respondWithError(w, http.StatusBadRequest, "Invalid tournament filetype provided", err)
		return
	}

	//m_type_parts := strings.Split(m_type, "/")
	key := make([]byte, 32)
	_, err = rand.Read(key)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to generate random filename", err)
		return
	}

	name := base64.RawURLEncoding.EncodeToString(key)
	file_name := fmt.Sprintf("%s.%s", name, "xlsx")
	file_path := filepath.Join("/tmp", file_name)
	file_ptr, err := os.Create(file_path)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to create location file for tournament file", err)
		return
	}
	defer file_ptr.Close()

	reader := io.Reader(file)
	writer := io.Writer(file_ptr)
	_, err = io.Copy(writer, reader)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to store tournament into destination", err)
		return
	}

	type parameters struct {
		Sheets       string `json:"excelSheets"`
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

	sheets_list := sheetsStringToSlice(params.Sheets)

	bcs_args := bsc.ExecutionArguments{
		Command:      "insert",
		DBName:       filepath.Join(api_config.db_dir, league.DatabaseName),
		ExcelFile:    file_path,
		ExcelSheets:  sheets_list,
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
		Name        string `json:"name"`
		Version     string `json:"version"`
		Tournaments []struct {
			Message string `json:"message"`
			Status  string `json:"status"`
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

func (api_config *apiConfig) handlerGetAllTournamentsInLeague(w http.ResponseWriter, r *http.Request) {
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

func (api_config *apiConfig) handlerGetTournamentResults(w http.ResponseWriter, r *http.Request) {
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

	tournament_id_string := r.PathValue("tournamentID")
	/*tournament_id, err := strconv.Atoi(tournament_id_string)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid tournament ID", err)
		return
	}*/

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

func sheetsStringToSlice(str string) []string {
	trimmed_str := strings.TrimRight(str, ";")
	var result []string
	if strings.Contains(trimmed_str, ";") {
		result = strings.Split(trimmed_str, ";")
	} else {
		result = append(result, trimmed_str)
	}
	return result
}
