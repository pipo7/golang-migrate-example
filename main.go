package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	db, err := sql.Open("postgres", "postgres://testuser:newpassword@0.0.0.0:5432/testdb?sslmode=disable")
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Could not create database driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		"postgres", driver)

	if err != nil {
		log.Fatalf("Could not create migrate instance: %v", err)
	}

	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration failed: %v", err)
	} else {
		fmt.Println("Down migrations were done sucessfully ")
	}

	// Running only certain steps of version
	_ = m.Steps(2) // NOTE This executes only 2 steps that is versions 1000 and 2000 thus After this
	// statement you see the session_review column is NOT added

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration failed: %v", err)
	} else {
		fmt.Println("Up migrations were done sucessfully ")
	}
}
