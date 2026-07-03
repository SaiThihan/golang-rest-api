package store

import (
	"database/sql"
	"fmt"
	"io/fs"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func Open() (*sql.DB, error) {
	connectionString := os.Getenv("DATABASE_URL")

	if connectionString == "" {
		return nil, fmt.Errorf("db: connection string is empty")
	}

	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		return nil, fmt.Errorf("db: open error %w", err)

	}

	fmt.Println("Database connected")
	return db, nil
}

func MigrateFs(db *sql.DB, migrationFs fs.FS, dir string) error {

	goose.SetBaseFS(migrationFs)

	defer func() {
		goose.SetBaseFS(nil)
	}()

	return Migrate(db, dir)
}

func Migrate(db *sql.DB, dir string) error {
	err := goose.SetDialect("postgres")

	if err != nil {
		return fmt.Errorf("db: set dialect error %w", err)
	}

	err = goose.Up(db, dir)
	if err != nil {
		return fmt.Errorf("db: migrate error %w", err)
	}
	return nil
}
