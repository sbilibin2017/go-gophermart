package tests

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/go-resty/resty/v2"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupContainer() testcontainers.Container {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpassword",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(60 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("failed to start container: %v", err)
	}

	return container
}

func setupDB(container testcontainers.Container) (*sql.DB, error) {
	ctx := context.Background()

	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return nil, fmt.Errorf("failed to get mapped port: %v", err)
	}

	databaseURL := fmt.Sprintf("postgres://testuser:testpassword@localhost:%s/testdb?sslmode=disable", port.Port())

	// Establish DB connection
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %v", err)
	}

	// Wait for DB to be ready
	time.Sleep(3 * time.Second)
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return db, nil
}

func runMigrations(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS rewards (
		match VARCHAR(255) PRIMARY KEY,
		reward BIGINT NOT NULL,
		reward_type VARCHAR(50) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW()
	);`

	_, err := db.Exec(query)
	return err
}

func clearTable(db *sql.DB) {
	_, err := db.Exec("DELETE FROM rewards")
	if err != nil {
		log.Fatalf("failed to clear rewards table: %v", err)
	}
}

func startServer(container testcontainers.Container) *exec.Cmd {
	ctx := context.Background()

	// Получаем порт и хост контейнера
	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		log.Fatalf("failed to get mapped port: %v", err)
	}
	host, err := container.Host(ctx)
	if err != nil {
		log.Fatalf("failed to get container host: %v", err)
	}

	// Строим корректный DATABASE_URI для контейнера
	dsn := fmt.Sprintf("postgres://testuser:testpassword@%s:%s/testdb?sslmode=disable", host, port.Port())

	// Формируем команду с флагом -d
	cmd := exec.Command("go", "run", "../main.go", "-d", dsn)

	// Пробрасываем вывод сервера в консоль
	cmd.Stdout = log.Writer()
	cmd.Stderr = log.Writer()

	// Запускаем
	err = cmd.Start()
	if err != nil {
		log.Fatalf("failed to start the server: %v", err)
	}

	// Дожидаемся готовности
	waitForServer("http://localhost:8081/api/goods", 10*time.Second)

	return cmd
}

func waitForServer(url string, timeout time.Duration) {
	client := resty.New()
	start := time.Now()
	for time.Since(start) < timeout {
		_, err := client.R().Get(url)
		if err == nil {
			return
		}
		time.Sleep(time.Second)
	}
	log.Fatalf("server did not start within %v", timeout)
}

func stopServer(cmd *exec.Cmd) {
	err := cmd.Process.Kill()
	if err != nil {
		log.Fatalf("failed to stop the server: %v", err)
	}
}
