package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    _ "github.com/lib/pq" // Import the PostgreSQL driver
    "log"
    "net/http"
)

type User struct {
    ID       int    `json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

var db *sql.DB

func main() {
    // Connection string to connect to PostgreSQL database
    connStr := "user=postgres password=abc@123 dbname=mydb sslmode=disable"
    
    // Open a connection to the database
    var err error
    db, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal("Failed to connect to the database:", err)
    }
    defer db.Close()

    // Verify the connection
    err = db.Ping()
    if err != nil {
        log.Fatal("Could not ping the database:", err)
    }

    fmt.Println("Connected to PostgreSQL!")

    // Create the users table if it doesn't exist
    createTableQuery := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        name TEXT NOT NULL,
        email TEXT NOT NULL,
        password TEXT NOT NULL
    );`
    _, err = db.Exec(createTableQuery)
    if err != nil {
        log.Fatal("Failed to create table:", err)
    }

    // Set up the HTTP server
    http.HandleFunc("/users", createUserHandler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}

// Handler function to create a new user
func createUserHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    var user User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Insert data into the database
    _, err = db.Exec("INSERT INTO users (id, name, email, password) VALUES ($1, $2, $3, $4)", user.ID, user.Name, user.Email, user.Password)
    if err != nil {
        log.Printf("Failed to insert data: %v", err)
        http.Error(w, "Failed to insert data", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    fmt.Fprintf(w, "User created successfully!")
}