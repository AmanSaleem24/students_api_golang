package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/AmanSaleem24/students-api/internal/config"
	"github.com/AmanSaleem24/students-api/internal/types"
	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	Db *sql.DB
}

func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {
	stat, err := s.Db.Prepare("INSERT INTO students (name, email, age) VALUES (?, ?, ?)")

	if err != nil {
		return 0, err
	}

	defer stat.Close()

	result, e := stat.Exec(name, email, age)
	if e != nil {
		return 0, e
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastId, nil

}

func (s *Sqlite) GetStudentById(id int64) (types.Student, error) {
	stat, err := s.Db.Prepare("SELECT * FROM students WHERE id = ? LIMIT 1")
	if err != nil {
		return types.Student{}, err
	}

	defer stat.Close()

	var student types.Student

	e := stat.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)
	if e != nil {
		if errors.Is(e, sql.ErrNoRows) {
			return types.Student{}, fmt.Errorf("no student found with id %d", id)
		}
		return types.Student{}, fmt.Errorf("Query error: %w", e)
	}

	return student, nil
}

func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	email TEXT,
	age INTEGER
	)`)

	if err != nil {
		return nil, err
	}

	return &Sqlite{
		Db: db,
	}, nil
}
