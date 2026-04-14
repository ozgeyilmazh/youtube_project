package quotesdiscovery

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/joho/godotenv"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	testDB         *sql.DB
	testTableName  string
	testRepository *ClickhouseRepository
)

//go:embed test_schema.sql
var testSchemaSQLTemplate string

func TestMain(m *testing.M) {
	ctx := context.Background()

	var err error
	testDB, testTableName, err = setupTestDB(ctx)
	if err != nil {
		fmt.Printf("Failed to setup test DB: %v", err)
		os.Exit(1)
	}
	testRepository = NewClickhouseRepository(testDB, testTableName)

	code := m.Run()

	_ = teardownTestDatabase(ctx, testDB, testTableName)

	os.Exit(code)
}

func setupTestDB(ctx context.Context) (*sql.DB, string, error) {
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

	clickhouseDatabase := os.Getenv("CLICKHOUSE_TEST_DATABASE")
	if clickhouseDatabase == "" {
		return nil, "", fmt.Errorf("CLICKHOUSE_DATABASE is not set")
	}

	clickhouseUsername := os.Getenv("CLICKHOUSE_USERNAME")
	if clickhouseUsername == "" {
		clickhouseUsername = "default"
	}
	clickhousePassword := os.Getenv("CLICKHOUSE_PASSWORD")
	tableName := os.Getenv("CLICKHOUSE_TEST_TABLE_NAME")
	if tableName == "" {
		tableName = "quotes_test"
	}

	db, err := sql.Open("clickhouse", fmt.Sprintf("clickhouse://%s:%s@%s:%s/%s", clickhouseUsername, clickhousePassword, clickhouseHost, clickhousePort, clickhouseDatabase))
	if err != nil {
		return nil, "", fmt.Errorf("Error opening clickhouse connection: %v", err)
	}

	err = db.PingContext(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("Error pinging clickhouse: %v", err)
	}

	createdDbQuery := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", clickhouseDatabase)
	_, err = db.ExecContext(ctx, createdDbQuery)
	if err != nil {
		return nil, "", fmt.Errorf("Error creating database: %v", err)
	}

	createTableSQL := fmt.Sprintf(strings.TrimSpace(testSchemaSQLTemplate), tableName)
	_, err = db.ExecContext(ctx, createTableSQL)
	if err != nil {
		_ = db.Close()
		return nil, "", fmt.Errorf("Error creating table: %v", err)
	}

	return db, tableName, nil
}

func teardownTestDatabase(ctx context.Context, db *sql.DB, tableName string) error {
	_, err := db.ExecContext(ctx, fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName))
	if err != nil {
		return fmt.Errorf("Error dropping table: %v", err)
	}
	return nil
}

func TestBulkInsertData(t *testing.T) {
	ctx := context.Background()
	Convey("Given a ClickhouseRepository is provided", t, func() {
		testData := []Quotes{
			{
				Quote:  "Test Quote",
				Author: "Test Author",
			},
			{
				Quote:  "Test Quote 2",
				Author: "Test Author 2",
			},
		}
		Convey("When BulkInsertData is called", func() {
			err := testRepository.BulkInsertData(ctx, testData)
			Convey("Then the error is nil", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When RunQuery is called", func() {
			query := fmt.Sprintf("SELECT * FROM %s", testTableName)
			err := testRepository.RunQuery(ctx, query)
			Convey("Then the error is nil", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}
