package integration_test

import (
	"database/sql"
	"flag"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var integration = flag.Bool("integration", false, "run integration tests")

func TestDBConnection(t *testing.T) {
	if !*integration {
		t.Skip("skipping integration test (use -integration flag to enable)")
	}

	dsn := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal("failed to connect to db:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Fatal("db is not reachable:", err)
	}

	t.Log("db connection successful")
}
