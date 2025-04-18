package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

type Repository struct {
	db *sql.DB
}

func NewRepository() (*Repository, error) {
	var dbPath string

	if runtime.GOOS == "windows" {
		appDataPath := os.Getenv("APPDATA")
		if appDataPath == "" {
			return nil, errors.New("the APPDATA environment variable is not set")
		}
		dbPath = filepath.Join(appDataPath, "gplan", "data.db")
	} else {
		homePath := os.Getenv("HOME")
		if homePath == "" {
			return nil, errors.New("the HOME environment variable is not set")
		}
		dbPath = filepath.Join(homePath, ".local", "share", "gplan", "data.db")
	}

	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0o755); err != nil {
		return nil, fmt.Errorf("error creating database directory: %w", err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	repo := &Repository{db: db}
	err = repo.init()
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *Repository) init() error {
	createProjectTable := `
	CREATE TABLE IF NOT EXISTS Projects (
		ID INTEGER PRIMARY KEY AUTOINCREMENT,
		Name TEXT NOT NULL UNIQUE,
		Description TEXT,
		CreationDate DATETIME NOT NULL,
		LastModifiedDate DATETIME NOT NULL
	);`

	createTaskTable := `
	CREATE TABLE IF NOT EXISTS Tasks (
		ID INTEGER PRIMARY KEY AUTOINCREMENT,
		Name TEXT NOT NULL,
		Description TEXT,
		ProjectID INTEGER NOT NULL,
		TaskCompleted BOOLEAN NOT NULL,
		DueDate DATETIME,
		CompletionDate DATETIME,
		CreationDate DATETIME NOT NULL,
		LastUpdatedDate DATETIME NOT NULL,
		Priority INTEGER NOT NULL DEFAULT 4,
		FOREIGN KEY (ProjectID) REFERENCES Projects(ID)
	);`

	_, err := r.db.Exec(createProjectTable)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(createTaskTable)
	return err
}

func (r *Repository) Close() {
	r.db.Close()
}
