package database

import (
	"database/sql"
	"fmt"
	"log"
)

type Database struct {
	Db *sql.DB
}

func NewDatabase() (*Database, error) {
	db, err := sql.Open("", "")
	if err != nil {
		return nil, err
	}

	// Check if the database is accessible
	err = db.Ping()
	if err != nil {
		err := db.Close()
		if err != nil {
			return nil, err
		} // Ensure the connection is closed if Ping fails
		return nil, err
	}
	fmt.Println("Database connection established successfully.")
	return &Database{Db: db}, nil
}

func (d *Database) Close() {
	err := d.Db.Close()
	if err != nil {
		log.Println("Error closing database connection:", err)
		return
	}
}

func (d *Database) GetDB() *sql.DB {
	return d.Db
}

func InitDatabase() (*Database, error) {
	_, err := NewDatabase()
	if err != nil {
		log.Println("Error initializing database:", err)
		return nil, err
	}
	fmt.Println("Database initialized successfully.")
	return nil, nil
}
