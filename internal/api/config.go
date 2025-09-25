package api

import (
	"database/sql"
	"log"
	"os"
	"sync/atomic"

	"github.com/Johnermac/http-server/internal/database"
	"github.com/joho/godotenv"
)

type APIConfig struct {
	FileserverHits 	atomic.Int32
	DB   					 *database.Queries
	Platform 				string
	JWTSecret  			string
	Polka_KEY				string
}


func newDB() *database.Queries {
    godotenv.Load()

    dbURL := os.Getenv("DB_URL")			

    db, err := sql.Open("postgres", dbURL)
    if err != nil {
        log.Fatal("cannot connect to db:", err)
    }

    return database.New(db)
}

func NewAPIConfig() *APIConfig {
	return &APIConfig{
		DB: newDB(),	
		Platform: os.Getenv("PLATFORM"),	
		JWTSecret: os.Getenv("JWT_SECRET"),
		Polka_KEY: os.Getenv("POLKA_KEY"),
	}
}
