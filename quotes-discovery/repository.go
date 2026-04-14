package quotesdiscovery

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"

	_ "github.com/ClickHouse/clickhouse-go/v2"
)

type Repository interface {
	BulkInsertData(ctx context.Context, quotes []Quotes) error
}

type ClickhouseRepository struct {
	db        *sql.DB
	tableName string
}

func SetupClickhouseRepository(tableName string) (*sql.DB, string, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, "", fmt.Errorf("Error loading .env file: %v", err)
	}

	clickhouseHost := os.Getenv("CLICKHOUSE_HOST")
	if clickhouseHost == "" {
		return nil, "", fmt.Errorf("CLICKHOUSE_HOST is not set")
	}
	clickhousePort := os.Getenv("CLICKHOUSE_PORT")
	if clickhousePort == "" {
		clickhousePort = "9000"
	}

	clickhouseDatabase := os.Getenv("CLICKHOUSE_DATABASE")
	clickhouseUsername := os.Getenv("CLICKHOUSE_USERNAME")
	if clickhouseUsername == "" {
		clickhouseUsername = "default"
	}
	clickhousePassword := os.Getenv("CLICKHOUSE_PASSWORD")

	tableName = os.Getenv("CLICKHOUSE_TABLE_NAME")
	if tableName == "" {
		tableName = "quotes"
	}

	db, err := sql.Open("clickhouse", fmt.Sprintf("clickhouse://%s:%s@%s:%s/%s", clickhouseUsername, clickhousePassword, clickhouseHost, clickhousePort, clickhouseDatabase))
	if err != nil {
		return nil, "", fmt.Errorf("Error opening clickhouse connection: %v", err)
	}
	return db, tableName, nil
}

func NewClickhouseRepository(db *sql.DB, tableName string) *ClickhouseRepository {
	return &ClickhouseRepository{db: db, tableName: tableName}
}

func (r *ClickhouseRepository) BulkInsertData(ctx context.Context, quotes []Quotes) error {
	if len(quotes) == 0 {
		return fmt.Errorf("No quotes to insert")
	}

	ts := time.Now().UTC().Truncate(time.Second)
	baseQuery := fmt.Sprintf(`INSERT INTO %s (id, quote, author, created_at, updated_at) VALUES`, r.tableName)
	valueStrings := make([]string, 0, len(quotes))
	valueArgs := make([]interface{}, 0, len(quotes)*5)
	for _, quote := range quotes {
		id := uuid.New().String()
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?)")
		valueArgs = append(valueArgs, id, quote.Quote, quote.Author, ts, ts)
	}
	query := baseQuery + strings.Join(valueStrings, ",")

	_, err := r.db.ExecContext(ctx, query, valueArgs...)
	if err != nil {
		return fmt.Errorf("Error inserting data: %v", err)
	}
	return nil
}

func (r *ClickhouseRepository) RunQuery(ctx context.Context, query string) error {
	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("Error running query: %v", err)
	}
	return nil
}
