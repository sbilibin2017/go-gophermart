package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type GoodsHandlerTestSuite struct {
	suite.Suite
	db            *sqlx.DB
	terminate     func()
	serverAddress string
}

func TestGoodsHandler(t *testing.T) {
	suite.Run(t, new(GoodsHandlerTestSuite))
}

func (suite *GoodsHandlerTestSuite) SetupSuite() {
	log.Println("Starting SetupSuite...")

	dsn, terminate, err := setupPostgresContainer(context.Background())
	suite.Require().NoError(err)

	log.Printf("PostgreSQL container started with DSN: %s\n", dsn)

	db, err := sqlx.Connect("pgx", dsn)
	suite.Require().NoError(err)
	suite.db = db
	suite.terminate = terminate
	setupDB(db)

	suite.serverAddress = "http://localhost:8080"

	runAddress = "localhost:8080"
	databaseURI = dsn

	log.Printf("Configuration - RunAddress: %s, DatabaseURI: %s", runAddress, databaseURI)

	go run()

	log.Println("Waiting for the server to start...")
	time.Sleep(5 * time.Second)
	log.Println("Server is ready.")
}

func (suite *GoodsHandlerTestSuite) TearDownSuite() {
	log.Println("Tearing down the test suite...")
	if suite.terminate != nil {
		suite.terminate()
	}
	if suite.db != nil {
		_ = suite.db.Close()
	}
	log.Println("Test suite teardown complete.")
}

func (suite *GoodsHandlerTestSuite) BeforeTest(suiteName, testName string) {
	log.Printf("Before test %s/%s: Clearing the database...\n", suiteName, testName)

	_, err := suite.db.Exec("TRUNCATE TABLE rewards RESTART IDENTITY CASCADE")
	if err != nil {
		log.Fatalf("Failed to truncate table: %v", err)
	}
	log.Println("Database cleared.")
}

func setupPostgresContainer(ctx context.Context) (dsn string, terminate func(), err error) {
	log.Println("Setting up PostgreSQL container...")

	req := testcontainers.ContainerRequest{
		Image:        "postgres:15",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": "password",
			"POSTGRES_USER":     "user",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").
			WithPollInterval(2 * time.Second).
			WithStartupTimeout(30 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return "", nil, fmt.Errorf("failed to start container: %w", err)
	}

	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "5432")
	dsn = fmt.Sprintf("postgres://user:password@%s:%s/testdb?sslmode=disable", host, port.Port())

	log.Printf("PostgreSQL container started at %s:%s\n", host, port.Port())

	terminate = func() {
		log.Println("Terminating PostgreSQL container...")
		_ = container.Terminate(ctx)
		log.Println("PostgreSQL container terminated.")
	}

	return dsn, terminate, nil
}

func setupDB(db *sqlx.DB) error {
	log.Println("Setting up database schema...")

	schema := `
	CREATE TABLE IF NOT EXISTS rewards (
		match VARCHAR(255) PRIMARY KEY,
		reward BIGINT NOT NULL,
		reward_type VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW()
	);`
	_, err := db.Exec(schema)
	if err != nil {
		log.Printf("Error creating schema: %v\n", err)
		return fmt.Errorf("failed to create schema: %w", err)
	}

	log.Println("Database schema set up successfully.")
	return nil
}

func (suite *GoodsHandlerTestSuite) TestRegisterReward() {
	tests := []struct {
		name         string
		reqBody      []map[string]interface{}
		expectedCode int
	}{
		{
			name: "Valid registration of reward",
			reqBody: []map[string]interface{}{
				{
					"match":       "Bork",
					"reward":      10,
					"reward_type": "%",
				},
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Empty match value",
			reqBody: []map[string]interface{}{
				{
					"match":       "",
					"reward":      10,
					"reward_type": "%",
				},
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Invalid reward value (non-positive reward)",
			reqBody: []map[string]interface{}{
				{
					"match":       "Bork",
					"reward":      0,
					"reward_type": "%",
				},
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Invalid reward type",
			reqBody: []map[string]interface{}{
				{
					"match":       "Bork",
					"reward":      10,
					"reward_type": "invalid_type",
				},
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Invalid request format (unknown field)",
			reqBody: []map[string]interface{}{
				{
					"unknown_field": "some_value",
				},
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Duplicate match registration attempt",
			reqBody: []map[string]interface{}{
				{
					"match":       "Bork",
					"reward":      10,
					"reward_type": "%",
				},
				{
					"match":       "Bork",
					"reward":      10,
					"reward_type": "%",
				},
			},
			expectedCode: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			for _, req := range tt.reqBody {
				reqData, err := json.Marshal(req)
				suite.Require().NoError(err)

				client := resty.New()
				resp, err := client.R().
					SetHeader("Content-Type", "application/json").
					SetBody(reqData).
					Post(suite.serverAddress + "/api/goods")

				log.Printf("Test: %s, Response status: %d\n", tt.name, resp.StatusCode())
				log.Printf("Response body: %s\n", resp.String())

				suite.Require().NoError(err)
				assert.Equal(suite.T(), tt.expectedCode, resp.StatusCode())
			}
		})
	}
}
