package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mahmoud-abd-elnasser/todo_api/internal/models"
)

func CreateTodo(pool *pgxpool.Pool, title string, completed bool, userID string) (*models.Todo, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	 INSERT INTO todos (title, completed, user_id)
	 VALUES ($1, $2, $3)
	 RETURNING id, title, completed, created_at, updated_at, user_id
	 `

	var todo models.Todo

	err := pool.QueryRow(ctx, query, title, completed, userID).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
		&todo.UserID,
	)

	if err != nil {
		return nil, err
	}

	return &todo, nil

}

func UpdateTodo(pool *pgxpool.Pool, id int, title string, completed bool, userID string) (*models.Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	UPDATE todos
	SET title = $1, completed = $2, updated_at = NOW()
	WHERE id = $3 AND user_id = $4
	RETURNING id, title, completed, created_at, updated_at, user_id
	`

	var todo models.Todo
	err := pool.QueryRow(ctx, query, title, completed, id, userID).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
		&todo.UserID,
	)
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func GetAllTodos(pool *pgxpool.Pool) ([]*models.Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	SELECT id, title, completed, created_at, updated_at, user_id
	FROM todos
	ORDER BY id ASC
	`
	todos := []*models.Todo{}

	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var todo models.Todo
		if err := rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.Completed,
			&todo.CreatedAt,
			&todo.UpdatedAt,
			&todo.UserID,
		); err != nil {
			return nil, err
		}
		todos = append(todos, &todo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

func GetTodoById(pool *pgxpool.Pool, id int) (*models.Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	SELECT id, title, completed, created_at, updated_at, user_id
    FROM todos 
    WHERE id = $1
	`
	var todo models.Todo

	err := pool.QueryRow(ctx, query, id).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
		&todo.UserID,
	)
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func DeleteTodo(pool *pgxpool.Pool, id int, userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	query := `
	DELETE FROM todos
	WHERE id = $1 AND user_id = $2
	`
	commandTag, err := pool.Exec(ctx, query, id, userID)
    if err != nil {
        return err
	}

    if commandTag.RowsAffected() == 0 {
        return errors.New("todo not found or unauthorized") 
    }

    return nil

}
