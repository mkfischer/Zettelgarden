package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go-backend/models"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Server struct {
	db      *sql.DB
	s3      *s3.Client
	testing bool
}

var s *Server

func helloWorld(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, "hello world")
}

func (s *Server) getAllFiles(w http.ResponseWriter, r *http.Request) {
	rows, _ := s.db.Query(`
	SELECT files.id, files.name, files.type, files.path, files.filename, files.size, files.created_by, files.updated_by, files.card_pk, files.created_at, files.updated_at
	FROM files
	JOIN cards ON files.card_pk = cards.id
	WHERE files.is_deleted = FALSE AND cards.user_id = $1`, 1)

	defer rows.Close()

	var files []models.File

	for rows.Next() {
		var file models.File
		if err := rows.Scan(
			&file.ID,
			&file.Name,
			&file.Filetype,
			&file.Path,
			&file.Filename,
			&file.Size,
			&file.CreatedBy,
			&file.UpdatedBy,
			&file.CardPK,
			&file.CreatedAt,
			&file.UpdatedAt,
		); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		files = append(files, file)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsonResponse, err := json.Marshal(files)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the JSON response
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func main() {
	s = &Server{}

	dbConfig := databaseConfig{}
	dbConfig.host = os.Getenv("DB_HOST")
	dbConfig.port = os.Getenv("DB_PORT")
	dbConfig.user = os.Getenv("DB_USER")
	dbConfig.password = os.Getenv("DB_PASS")
	dbConfig.databaseName = os.Getenv("DB_NAME")

	db, err := ConnectToDatabase(dbConfig)
	if err != nil {
		log.Fatalf("Unable to connect to the database: %v\n", err)
	}
	s.db = db
	s.s3 = createS3Client()

	//http.HandleFunc("GET /api/files", getAllFiles)
	//http.HandleFunc("POST /api/files/upload", uplpadFile)
	//http.HandleFunc("GET /api/files/{I}/", getFileMetadata)
	//http.HandleFunc("PATCH /api/files/{I}/", editFile)
	//http.HandleFunc("DELETE /api/files/{I}/", deleteFile)
	//http.HandleFunc("GET /api/files/download/{id}", helloWorld)
	http.ListenAndServe(":8080", nil)
}
