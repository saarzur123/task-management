package service

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"github.com/saarzur123/task-management/backend/models"
	"strconv"
	"time"
)

type TaskRepository interface {
	Create(task *models.Task) error
	GetByID(id string) (*models.Task, error)
	Update(task *models.Task) error
	Delete(id string) error
	GetAll() ([]models.Task, error)
}

type TaskManager struct {
	DB *sql.DB
}

var (
	ErrNotFound = errors.New("NotFound")
)

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		return nil, err
	}

	sqlStmt := "CREATE TABLE IF NOT EXISTS tasks (" +
		"id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT," +
		"title TEXT NOT NULL," +
		"description TEXT NOT NULL," +
		"status TEXT NOT NULL," +
		"created_at TIMESTAMP NOT NULL);"

	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (m *TaskManager) Create(task *models.Task) error {
	task.CreatedAt = time.Now()
	query := `INSERT INTO tasks (title, description, status, created_at) VALUES (?, ?, ?, ?)`
	row, err := m.DB.Exec(query, task.Title, task.Description, task.Status, task.CreatedAt)
	if err != nil {
		return err
	}

	dbId, err := row.LastInsertId()
	if err != nil {
		return err
	}

	task.Id = strconv.FormatInt(dbId, 10)
	return nil
}

func (m *TaskManager) GetByID(id string) (*models.Task, error) {
	query := `SELECT id, title, description, status, created_at FROM tasks WHERE id = ?`
	row := m.DB.QueryRow(query, id)

	task := models.Task{}
	err := row.Scan(&task.Id, &task.Title, &task.Description, &task.Status, &task.CreatedAt)
	if err != nil {
		return &task, err
	}

	return &task, nil
}

func (m *TaskManager) Update(task *models.Task) error {
	query := `UPDATE tasks SET title = ?, description = ?, status = ? WHERE id = ?`
	rows, err := m.DB.Exec(query, task.Title, task.Description, task.Status, task.Id)
	if err != nil {
		return err
	}

	rowsAffected, err := rows.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}
	return err
}

func (m *TaskManager) Delete(id string) error {
	query := `DELETE FROM tasks WHERE id = ?`
	rows, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := rows.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return err
}

func (m *TaskManager) GetAll() ([]models.Task, error) {
	query := `SELECT id, title, description, status, created_at FROM tasks`
	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]models.Task, 0)
	for rows.Next() {
		var task models.Task
		err = rows.Scan(&task.Id, &task.Title, &task.Description, &task.Status, &task.CreatedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
