package repository

import (
	"database/sql"
	"log"
	"todo/model"
)

func GetTodos(db *sql.DB) ([]model.Todo, error) {
	rows, err := db.Query("SELECT id, title, completed FROM todos ORDER BY id")

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var todos []model.Todo
	for rows.Next() {
		var todo model.Todo
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Completed)

		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

func CreateTodo(db *sql.DB, todo model.Todo) (model.Todo, error) {
	var id int

	err := db.QueryRow("INSERT INTO todos (title, completed) VALUES ($1, $2) RETURNING id", todo.Title, todo.Completed).Scan(&id)

	if err != nil {
		return model.Todo{}, err
	}

	todo.ID = id
	return todo, nil
}

func DeleteTodo(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM todos WHERE id = $1", id)

	if err != nil {
		return err
	}

	return nil
}

func UpdateTodo(db *sql.DB, todo model.Todo) (model.Todo, error) {
	_, err := db.Exec("UPDATE todos SET title  = $1, completed = $2 WHERE id = $3", todo.Title, todo.Completed, todo.ID)

	if err != nil {
		return model.Todo{}, err
	}

	return todo, nil
}
