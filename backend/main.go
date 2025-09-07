package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"
	"todo/handler"
	"todo/middleware"

	_ "github.com/lib/pq"
)

func createTableIfNotExists(db *sql.DB) {
	query := `
		CREATE TABLE IF NOT EXISTS todos (
			id BIGSERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			completed BOOLEAN DEFAULT FALSE
		);`

	_, err := db.Exec(query)

	if err != nil {
		log.Fatal(err)
	}

	println("Table created successfully or already exists")
}

func main() {
	mux := http.NewServeMux()

	db, err := sql.Open("postgres", "user=admin password=admin dbname=todo sslmode=disable")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	mux.HandleFunc("GET /api/todos", handler.GetTodos(db))
	mux.HandleFunc("POST /api/todo", handler.CreateTodo(db))
	mux.HandleFunc("DELETE /api/todo/{id}", handler.DeleteTodo(db))
	mux.HandleFunc("PUT /api/todo/{id}", handler.UpdateTodo(db))

	println("Connected to the database successfully")
	createTableIfNotExists(db)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      middleware.CorsMiddleware(mux),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("Starting server on :8080")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
