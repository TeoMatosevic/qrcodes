package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

type Service interface {
	Close() error

	CreateIfNotExists() error

	TableName() string

	Insert(t Ticket) error

	Count() (int, error)

	AmountByVatin(vatin string) (int, error)
}

type service struct {
	db         *sql.DB
	table_name string
}

var (
	tableName  = os.Getenv("QRCODES_DB_TABLE_NAME")
	connStr    = os.Getenv("QRCODES_DB_CONNECTION_STRING")
	dbInstance *service
)

func New() Service {
	if dbInstance != nil {
		return dbInstance
	}

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}

	dbInstance = &service{
		db: db,
	}

	dbInstance.table_name = tableName
	err = dbInstance.CreateIfNotExists()
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}

	return dbInstance
}

func (s *service) Close() error {
	return s.db.Close()
}

func (s *service) CreateIfNotExists() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, fmt.Sprintf(`
        CREATE TABLE IF NOT EXISTS %s (
            id TEXT PRIMARY KEY,
            vatin TEXT NOT NULL,
            first_name TEXT NOT NULL,
            last_name TEXT NOT NULL,
            created_at TEXT NOT NULL
        )`, s.table_name))

	if err != nil {
		return err
	}

	return nil
}

func (s *service) TableName() string {
	return s.table_name
}

func (s *service) Insert(t Ticket) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, fmt.Sprintf(`
        INSERT INTO %s (id, vatin, first_name, last_name, created_at)
        VALUES ($1, $2, $3, $4, $5)`, s.table_name), t.ID, t.Vatin, t.FirstName, t.LastName, t.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (s *service) Count() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var count int
	err := s.db.QueryRowContext(ctx, fmt.Sprintf(`SELECT COUNT(*) FROM %s`, s.table_name)).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s *service) AmountByVatin(vatin string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var count int
	err := s.db.QueryRowContext(ctx, fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE vatin = $1`, s.table_name), vatin).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
