package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    _ "github.com/lib/pq" // Import the PostgreSQL driver
    "log"
    "net/http"
    "os"
    "strconv"
    "github.com/gorilla/mux"
    "github.com/joho/godotenv"
    "golang.org/x/crypto/bcrypt"
)

type User struct {
    ID       int    `json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

type Post struct {
    ID             int    `json:"id"`
    UserID         int    `json:"user_id"`
    Caption        string `json:"caption"`
    ImageURL       string `json:"image_url"`
    PostedTimestamp string `json:"posted_timestamp"`
}

var db *sql.DB

func main() {
    // Load environment variables from .env file
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    // Get the PostgreSQL credentials from environment variables
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")
    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")

    // Connection string to connect to PostgreSQL
    connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", dbUser, dbPassword, dbName, dbHost, dbPort)
    
    // Open a connection to the database
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
    createUsersTableQuery := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        name TEXT NOT NULL,
        email TEXT NOT NULL,
        password TEXT NOT NULL
    );`
    _, err = db.Exec(createUsersTableQuery)
    if err != nil {
        log.Fatal("Failed to create users table:", err)
    }

    // Create the posts table if it doesn't exist
    createPostsTableQuery := `
    CREATE TABLE IF NOT EXISTS posts (
        id SERIAL PRIMARY KEY,
        user_id INT NOT NULL,
        caption TEXT NOT NULL,
        image_url TEXT NOT NULL,
        posted_timestamp TIMESTAMP NOT NULL,
        FOREIGN KEY (user_id) REFERENCES users(id)
    );`
    _, err = db.Exec(createPostsTableQuery)
    if err != nil {
        log.Fatal("Failed to create posts table:", err)
    }

    // Set up the HTTP server
    r := mux.NewRouter()
    r.HandleFunc("/users", createUserHandler).Methods("POST")
    r.HandleFunc("/users/{id}", getUserHandler).Methods("GET")
    r.HandleFunc("/posts", createPostHandler).Methods("POST")
    r.HandleFunc("/posts/{id}", getPostHandler).Methods("GET")
    r.HandleFunc("/posts/users/{id}", getUserPostsHandler).Methods("GET")
    log.Fatal(http.ListenAndServe(":8080", r))
}

// Handler function to create a new user
func createUserHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    var user struct {
        Name     string `json:"name"`
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Hash the password using bcrypt
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        log.Printf("Failed to hash password: %v", err)
        http.Error(w, "Failed to hash password", http.StatusInternalServerError)
        return
    }

    // Insert data into the database, excluding the id field
    _, err = db.Exec("INSERT INTO users (name, email, password) VALUES ($1, $2, $3)", user.Name, user.Email, string(hashedPassword))
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

// Handler function to create a new post
func createPostHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    var post Post
    err := json.NewDecoder(r.Body).Decode(&post)
    if err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Insert data into the database, excluding the id field
    _, err = db.Exec("INSERT INTO posts (user_id, caption, image_url, posted_timestamp) VALUES ($1, $2, $3, $4)", post.UserID, post.Caption, post.ImageURL, post.PostedTimestamp)
    if err != nil {
        log.Printf("Failed to insert data: %v", err)
        http.Error(w, "Failed to insert data", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    fmt.Fprintf(w, "Post created successfully!")
}

// Handler function to get a post by ID
func getPostHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    idStr := vars["id"]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid post ID", http.StatusBadRequest)
        return
    }

    var post Post
    err = db.QueryRow("SELECT id, user_id, caption, image_url, posted_timestamp FROM posts WHERE id = $1", id).Scan(&post.ID, &post.UserID, &post.Caption, &post.ImageURL, &post.PostedTimestamp)
    if err != nil {
        if err == sql.ErrNoRows {
            http.Error(w, "Post not found", http.StatusNotFound)
        } else {
            log.Printf("Failed to query post: %v", err)
            http.Error(w, "Failed to query post", http.StatusInternalServerError)
        }
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(post)
}

// Handler function to get all posts by user ID
func getUserPostsHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    idStr := vars["id"]
    userID, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    rows, err := db.Query("SELECT id, user_id, caption, image_url, posted_timestamp FROM posts WHERE user_id = $1", userID)
    if err != nil {
        log.Printf("Failed to query posts: %v", err)
        http.Error(w, "Failed to query posts", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var posts []Post
    for rows.Next() {
        var post Post
        err := rows.Scan(&post.ID, &post.UserID, &post.Caption, &post.ImageURL, &post.PostedTimestamp)
        if err != nil {
            log.Printf("Failed to scan post: %v", err)
            http.Error(w, "Failed to scan post", http.StatusInternalServerError)
            return
        }
        posts = append(posts, post)
    }

    if err = rows.Err(); err != nil {
        log.Printf("Rows error: %v", err)
        http.Error(w, "Rows error", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(posts)
}