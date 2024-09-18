package main

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq" // Import the PostgreSQL driver
    "log"
    "net/http"
    "strconv"
)

// type User struct {
//     ID       int    `json:"id"`
//     Name     string `json:"name"`
//     Email    string `json:"email"`
//     Password string `json:"password"`
// }

var db *sql.DB

func main() {
    // Connection string to connect to the PostgreSQL database
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

    // Parse URL parameters
    idStr := r.URL.Query().Get("id")
    name := r.URL.Query().Get("name")
    email := r.URL.Query().Get("email")
    password := r.URL.Query().Get("password")

    if idStr == "" || name == "" || email == "" || password == "" {
        http.Error(w, "Missing URL parameters", http.StatusBadRequest)
        return
    }

    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid ID parameter", http.StatusBadRequest)
        return
    }

    // Insert data into the database
    _, err = db.Exec("INSERT INTO users (id, name, email, password) VALUES ($1, $2, $3, $4)", id, name, email, password)
    if err != nil {
        log.Printf("Failed to insert data: %v", err)
        http.Error(w, "Failed to insert data", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    fmt.Fprintf(w, "User created successfully!")
}