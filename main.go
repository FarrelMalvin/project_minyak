package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"project_minyak/routes"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("mysql", "root@tcp(127.0.0.1:3307)/project_minyak")
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	// Test koneksi ke database
	err = DB.Ping()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	log.Println("Database initialized successfully")
}

func main() {

	InitDB()

	r := routes.SetupRoutes(DB)

	port := ":9090"
	fmt.Println("Server running on port", port)
	log.Fatal(http.ListenAndServe(port, r))
}
