package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    _ "github.com/lib/pq" // Import the PostgreSQL driver
    "log"
    "net/http"
    "strconv"
    "github.com/gorilla/mux"
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
    r := mux.NewRouter()
    r.HandleFunc("/users", createUserHandler).Methods("POST")
    r.HandleFunc("/users/{id}", getUserHandler).Methods("GET")
    log.Fatal(http.ListenAndServe(":8080", r))
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

// Handler function to get a user by ID
func getUserHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    idStr := vars["id"]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    var user User
    err = db.QueryRow("SELECT id, name, email, password FROM users WHERE id = $1", id).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
    if err != nil {
        if err == sql.ErrNoRows {
            http.Error(w, "User not found", http.StatusNotFound)
        } else {
            log.Printf("Failed to query user: %v", err)
            http.Error(w, "Failed to query user", http.StatusInternalServerError)
        }
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}