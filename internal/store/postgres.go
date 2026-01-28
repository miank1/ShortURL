package store

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(connStr string) (*PostgresStore, error) {
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}

	// verify connection at startup
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) Save(code, url string, expiresAt time.Time) error {
	_, err := s.db.Exec(
		`INSERT INTO short_urls (code, long_url, expires_at)
		 VALUES ($1, $2, $3)`,
		code,
		url,
		expiresAt,
	)
	return err
}

func (s *PostgresStore) Get(code string) (string, time.Time, bool, error) {
	var longURL string
	var expiresAt time.Time

	err := s.db.QueryRow(
		`SELECT long_url, expires_at
		 FROM short_urls
		 WHERE code = $1`,
		code,
	).Scan(&longURL, &expiresAt)

	if err == sql.ErrNoRows {
		return "", time.Time{}, false, nil
	}
	if err != nil {
		return "", time.Time{}, false, err
	}

	return longURL, expiresAt, true, nil
}

func (s *PostgresStore) Close() error {
	return s.db.Close()
}
