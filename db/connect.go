package db

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	DB_HOSTNAME string
	DB_PORT     int
	DB_USERNAME string
	DB_PASSWORD string
	DB_NAME     string
)

const (
	FORUMPOST_URL = "https://osu.ppy.sh/community/forums/topics/"
	USER_URL      = "https://osu.ppy.sh/users/"
)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}
	dbHostname, exists := os.LookupEnv("POSTGRES_HOSTNAME")
	if exists {
		DB_HOSTNAME = dbHostname
	}
	dbUsername, exists := os.LookupEnv("POSTGRES_USER")
	if exists {
		DB_USERNAME = dbUsername
	}
	dbPassword, exists := os.LookupEnv("POSTGRES_PASSWORD")
	if exists {
		DB_PASSWORD = dbPassword
	}
	dbName, exists := os.LookupEnv("POSTGRES_DATABASE")
	if exists {
		DB_NAME = dbName
	}
	dbPort, exists := os.LookupEnv("POSTGRES_PORT")
	if exists {
		DB_PORT, _ = strconv.Atoi(dbPort)
	}
}

func ConnectDB() *sql.DB {
	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		DB_HOSTNAME, DB_PORT, DB_USERNAME, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println("Error connecting database:", err)
		return nil
	}
	return db
}
