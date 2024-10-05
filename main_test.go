package main

import (
	"os"
	"testing"
)

func TestSurrealDBConnection(t *testing.T) {
	dbURL := os.Getenv("SURREALDB_URL")
	dbNamespace := os.Getenv("SURREALDB_NAMESPACE")
	dbDatabase := os.Getenv("SURREALDB_DATABASE")

	if dbURL == "" || dbNamespace == "" || dbDatabase == "" {
		t.Fatalf("Environment variables not set correctly")
	}

	db := SurrealDB(Surreal{
		URL:       dbURL,
		Namespace: dbNamespace,
		Database:  dbDatabase,
		Version:   "1.x <=",
	})

	result, err := db.Query(`RETURN $message;`, map[string]interface{}{
		"message": "Hello World!",
	})
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	if len(result) == 0 {
		t.Fatalf("No results returned from the query")
	}

	if result[0].Status != "OK" {
		t.Fatalf("Query failed with errors: %v", result[0].Result)
	}

	result_, ok := result[0].Result.(string)
	if !ok {
		t.Fatalf("Unexpected Result")
	}

	t.Log("Status: ", result[0].Status)
	t.Log("Results: ", result_)
}
