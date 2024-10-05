# SurrealDB Go Client

This Go package provides an easy-to-use interface for connecting to and querying a SurrealDB database. It supports various authentication methods and can execute SurrealQL queries.

## Features

- Connect to a SurrealDB instance with configurable parameters.
- Support for multiple authentication methods: Root, Token, and Scope.
- Execute queries and handle responses in a structured way.
- Sign in and sign up users with authentication token management.

## Installation

To use this package in your Go project, you can simply copy the source files into your project directory or use Go modules.

```bash
go get github.com/brahimEmo/go-surrealdb
```

## Usage

### Configuration

Create a new instance of SurrealDB by providing the required configuration parameters:

```go
package main

import (
	"fmt"
	"github.com/brahimEmo/go-surrealdb"
)

func main() {
	dbConfig := Surreal{
		URL:       "http://localhost:8000",
		Version:   ">= 2.x",
		Namespace: "<your namespace>",
		Database:  "<your database>",
		Auth: &AuthConfig{
			Method: "Root",
			Vars: AuthVars{
				Username: "<your username>",
				Password: "<your password>",
			},
		},
	}

	db := SurrealDB(dbConfig)

	// Now you can use the db object to perform operations.
}
```

### Querying the Database

You can execute SurrealQL queries using the `Query` method:

```go
query := "SELECT * FROM your_table"
vars := map[string]interface{}{
	"param1": "value1",
	"param2": 42,
}

results, err := db.Query(query, vars)
if err != nil {
	fmt.Println("Error executing query:", err)
	return
}

fmt.Println("Query Results:", results)
```

### Signing In

To sign in and obtain an authentication token, use the `Signin` method:

```go
signinVars := SigninVars{
	NS:   &db.Namespace,
	DB:   &db.Database,
	User: ptrString("<your username>"),
	Pass: ptrString("<your password>"),
	AC:   ptrString("<your access control>"), // Optional, based on your needs
}

token, err := db.Signin(signinVars)
if err != nil {
	fmt.Println("Error signing in:", err)
	return
}

fmt.Println("Authentication Token:", token)
```

### Signing Up

For signing up new users, use the `Signup` method:

```go
signupVars := SignupVars{
	NS:   &db.Namespace,
	DB:   &db.Database,
	User: ptrString("<new username>"),
	Pass: ptrString("<new password>"),
	AC:   ptrString("<your access control>"), // Optional, based on your needs
}

token, err := db.Signup(signupVars)
if err != nil {
	fmt.Println("Error signing up:", err)
	return
}

fmt.Println("Authentication Token:", token)
```

### Authentication Management

You can update or remove the authentication method at any time:

```go
newAuth := &AuthConfig{
	Method: "Token",
	Token:  "your_token",
}

db.Authenticate(newAuth)

// To remove authentication
db.Auth = nil
  OR
db.Invalidate()
```

### Error Handling

The package provides structured error handling through the `RequestError` type, which includes error codes and descriptions. Always check for errors when executing queries or authentication requests.

## License

This package is licensed under the MIT License. See the LICENSE file for details.

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request with your changes.