package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// Handles the HTTP connection and querying of a SurrealDB database.
type Surreal struct {
	URL       string      // URL of the SurrealDB instance (e.g., http://localhost:8000).
	Version   string      // The version of SurrealDB being used. Should be set to "1.x <=" for version 1.x or ">= 2.x" for version 2.x. Defaults to ">= 2.x" if not specified.
	Namespace string      // The namespace within the SurrealDB instance to use for database operations.
	Database  string      // The specific database to connect to within the specified namespace.
	Auth      *AuthConfig // Configuration for authentication, including methods such as Root, Token, or Scope. This field is optional and can be nil if no authentication is required.
}

// Represents the authentication configuration, supporting Root, Token, and Scope methods.
type AuthConfig struct {
	Method string   // Authentication method: "Root", "Token", or "Scope".
	Vars   AuthVars // Variables required for authentication (username, password, namespace, etc.).
	Token  string   // Token for token-based authentication.
	Scope  string   // Scope for scope-based authentication.
}

// Represents the credentials for root-based authentication.
type AuthVars struct {
	Username  string // Username for root-based authentication.
	Password  string // Password for root-based authentication.
	Namespace string // Namespace (optional) for root-based authentication.
	Database  string // Database (optional) for root-based authentication.
}

// Rrepresents the result of a query executed against SurrealDB.
// It contains the query result, status of the operation, and the time it took to process.
type QueryResult struct {
	Result interface{} `json:"result"` // The result of the query, which can be of any type.
	Status string      `json:"status"` // The status of the query operation ("OK" or "ERR").
	Time   string      `json:"time"`   // The time taken to process the query, formatted as a string.
}

// Represents an error returned from SurrealDB.
// It contains detailed information about the error, including a code, a description, and additional information.
type RequestError struct {
	Code        int    `json:"code"`        // A numerical code representing the error type (e.g., 400).
	Details     string `json:"details"`     // A brief description of the error details.
	Description string `json:"description"` // A human-readable description of the error.
	Information string `json:"information"` // Additional information that may help in diagnosing the issue.
}

// Creates a new instance of SurrealDB with the provided configuration.
func SurrealDB(config Surreal) *Surreal {
	if config.Version == "" || (config.Version != ">= 2.x" && config.Version != "1.x <=") {
		config.Version = ">= 2.x" // Default version if not provided.
	}
	return &config
}

// Configures the HTTP request with appropriate headers based on SurrealDB version and authentication.
func (db *Surreal) setRequest(url, method string, body interface{}, isRaw bool) (*http.Request, error) {
	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	// Set headers based on SurrealDB version.
	if db.Version == "1.x <=" {
		headers["NS"] = db.Namespace
		headers["DB"] = db.Database
	} else {
		headers["surreal-ns"] = db.Namespace
		headers["surreal-db"] = db.Database
	}

	// Handle authentication.
	if db.Auth != nil {
		switch db.Auth.Method {
		case "Root":
			authStr := fmt.Sprintf("%s:%s", db.Auth.Vars.Username, db.Auth.Vars.Password)
			headers["Authorization"] = "Basic " + base64.StdEncoding.EncodeToString([]byte(authStr))

			if db.Auth.Vars.Namespace != "" {
				headers["NS"] = db.Auth.Vars.Namespace
			}
			if db.Auth.Vars.Database != "" {
				headers["DB"] = db.Auth.Vars.Database
			}
		case "Token":
			headers["Authorization"] = "Bearer " + db.Auth.Token
		case "Scope":
			headers["SC"] = db.Auth.Scope
			for key, value := range db.Auth.Vars.ToMap() {
				headers[key] = value
			}
		}
	}

	var reqBody []byte
	var err error

	// Determine whether to send raw text or JSON
	if isRaw {
		reqBody = []byte(body.(string)) // Assert that body is a string
	} else {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	// Create the HTTP request with the raw SQL query as the body.
	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	// Set headers in the request.
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return req, nil
}

// Sends the HTTP request and processes the response, returning the result or error.
func (db *Surreal) processRequest(req *http.Request) ([]QueryResult, *RequestError) {
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, &RequestError{
			Code:        500,
			Details:     "Request problems detected",
			Description: "There is a problem with your request. Refer to the documentation for further information.",
			Information: err.Error(),
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var result RequestError
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			return nil, &RequestError{
				Code:        500,
				Details:     "Request problems detected",
				Description: "There is a problem with your request. Refer to the documentation for further information.",
				Information: err.Error(),
			}
		}

		return nil, &result
	}

	// Slice of QueryResult with a generic type
	var result []QueryResult
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, &RequestError{
			Code:        500,
			Details:     "Request problems detected",
			Description: "There is a problem with your request. Refer to the documentation for further information.",
			Information: err.Error(),
		}
	}

	return result, nil
}

// Sends a SurrealQL query to the SurrealDB server, with optional query variables.
func (db *Surreal) Query(sql string, vars map[string]interface{}) ([]QueryResult, *RequestError) {
	queryURL := fmt.Sprintf("%s/sql", db.URL)
	queryVars := url.Values{}

	// Append query variables.
	for key, value := range vars {
		if value != nil {
			queryVars.Add(key, fmt.Sprintf("%v", value))
		}
	}

	// If queryVars are present, add them to the URL.
	if encodedVars := queryVars.Encode(); encodedVars != "" {
		queryURL = fmt.Sprintf("%s?%s", queryURL, encodedVars)
	}

	// Set the request with appropriate headers and body (SurrealQL query as body).
	req, err := db.setRequest(queryURL, "POST", sql, true)
	if err != nil {
		return nil, &RequestError{
			Code:        500,
			Details:     "Request problems detected",
			Description: "There is a problem with your request. Refer to the documentation for further information.",
			Information: err.Error(),
		}
	}

	// Process the request and return the response.
	return db.processRequest(req)
}

// Updates the authentication method for the SurrealDB connection.
func (db *Surreal) Authenticate(auth *AuthConfig) {
	db.Auth = auth
}

// Removes the authentication method from the SurrealDB connection.
func (db *Surreal) Invalidate() {
	db.Auth = nil
}

// Converts the AuthVars struct to a map for scope-based authentication.
func (a AuthVars) ToMap() map[string]string {
	return map[string]string{
		"username":  a.Username,
		"password":  a.Password,
		"namespace": a.Namespace,
		"database":  a.Database,
	}
}
