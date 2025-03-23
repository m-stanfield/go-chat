package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"go-chat-react/internal/database"
)

var (
	dburl      = os.Getenv("BLUEPRINT_DB_URL")
	dbInstance *database.DBService
)

func executeSQLFile(db *sql.DB, filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	_, err = db.Exec(string(data))
	if err != nil {
		return fmt.Errorf("failed to execute SQL: %w", err)
	}

	return nil
}

func NewInMemoryDB() *database.DBService {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
	}

	err = executeSQLFile(db, "../../schema.sql")
	if err != nil {
		log.Fatal(err)
	}
	err = executeSQLFile(db, "../../mockdata.sql")
	if err != nil {
		log.Fatal(err)
	}

	return database.New(db)
}

func NewDB() *database.DBService {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}

	db, err := sql.Open("sqlite3", dburl)
	if err != nil {
		// This will not be a connection error, but a DSN parse error or
		// another initialization error.
		log.Fatal(err)
	}

	return database.New(db)
}

type Server struct {
	port int

	db database.Service
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	fmt.Printf("opening on port %d", port)
	NewServer := &Server{
		port: port,

		db: NewDB(),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
