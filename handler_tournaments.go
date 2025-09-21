package main

import (
	"encoding/base64"
	"path/filepath"
	"net/http"
	"strconv"
	"mime"
	"fmt"
	"os"
	"io"
	"crypto/rand"
	"errors"

	"github.com/t6kke/skill-calculator/internal/auth"
	"github.com/t6kke/skill-calculator/internal/bsc"
)

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

	const maxMemory = 10 << 20 // 10 MB using bit shifting
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

	//TODO here goes the BSC usage logic
	type tempReply struct {
		BSC_reply string `json:"bsc_reply"`
	}

	bcs_args := bsc.ExecutionArguments{
		DBName:       filepath.Join(api_config.db_dir, league.DatabaseName),
		ExcelFile:    file_path,
		ExcelSheet:   "Sheet1",
		CategoryName: "TDC",
		CategoryDesc: "test doubles category",
	}
	exit_code, output_str := bcs_args.BSCExecution()
	if exit_code != 0 {
		error_message := fmt.Sprintf("exit code: %d", exit_code)
		respondWithError(w, http.StatusInternalServerError, "BSC execution failed", errors.New(error_message))
	}
	response := tempReply{
		BSC_reply: output_str,
	}

	respondWithJSON(w, http.StatusOK, response)
}
