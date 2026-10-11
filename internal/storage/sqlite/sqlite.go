package sqlite

import (
	"database/sql"

	"github.com/AmanSaleem24/students-api/internal/config"
	_"github.com/mattn/go-sqlite3"
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

func New(cfg *config.Config) (*Sqlite, error){
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