package librarian

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

type PostgresStorage struct {
	DB *sql.DB
}

func NewPostgresStorage() (*PostgresStorage, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		viper.GetString("postgres.host"), viper.GetInt("postgres.port"),
		viper.GetString("postgres.user"), viper.GetString("postgres.password"),
		viper.GetString("postgres.dbname"))

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return &PostgresStorage{DB: db}, nil
}

func (p *PostgresStorage) Store(ctx context.Context, key string, value []byte) error {
	var js json.RawMessage = value
	_, err := p.DB.ExecContext(ctx, "INSERT INTO librarian (key, value) VALUES ($1, $2) ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value", key, js)
	return err
}

func (p *PostgresStorage) Retrieve(ctx context.Context, key string) ([]byte, error) {
	var value []byte
	err := p.DB.QueryRowContext(ctx, "SELECT value FROM librarian WHERE key = $1", key).Scan(&value)
	if err != nil {
		return nil, err
	}
	return value, nil
}
