package go_surrealdb

import (
	"os"
	"testing"
)

func TestSurrealDBConnection(t *testing.T) {
	dbURL := os.Getenv("SURREALDB_URL")
	dbNamespace := os.Getenv("SURREALDB_NAMESPACE")
	dbDatabase := os.Getenv("SURREALDB_DATABASE")

	if dbURL == "" || dbNamespace == "" || dbDatabase == "" {
		t.Skip("Environment variables not set correctly")
	}

	db := SurrealDB(Surreal{
		URL:       dbURL,
		Namespace: dbNamespace,
		Database:  dbDatabase,
		Version:   "1.x <=",
	})

	result, err := db.Query("RETURN 1+1;", map[string]interface{}{})
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	t.Log("Results: ", result)
}
