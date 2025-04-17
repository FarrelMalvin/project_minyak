package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"project_minyak/routes"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

var DB *sql.DB

func InitDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, dbname)
	fmt.Println("Connecting to:", dsn)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	DB = db
	log.Println("Database initialized successfully")
}

func main() {
	InitDB()

	r := routes.SetupRoutes(DB)

	port := ":9090"
	fmt.Println("Server running on port", port)

	log.Fatal(http.ListenAndServe(port, r))
}
