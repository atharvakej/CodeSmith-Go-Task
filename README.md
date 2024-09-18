# Go Backend Project

This project is a simple backend application written in Go. It provides APIs for user management and post management. The project uses PostgreSQL as the database and Gorilla Mux for routing.

## Features

- **User Management**: Create and retrieve users.
- **Post Management**: Create and retrieve posts by users.
- **Password Encryption**: Passwords are encrypted using bcrypt before storing them in the database.
- **Environment Variables**: Database credentials and other configurations are managed using a `.env` file.

## Getting Started

### Prerequisites

- Go 1.16 or later
- PostgreSQL
- Git

### Installation

1. Clone the repository:

    ```sh
    git clone https://github.com/yourusername/go-backend.git
    cd go-backend
    ```

2. Install the dependencies:

    ```sh
    go get -u github.com/gorilla/mux
    go get -u github.com/lib/pq
    go get -u github.com/joho/godotenv
    go get -u golang.org/x/crypto/bcrypt
    go get -u github.com/stretchr/testify
    ```

3. Create a `.env` file in the root directory and add your PostgreSQL credentials:

    ```env
    DB_USER=your_db_user
    DB_PASSWORD=your_db_password
    DB_NAME=your_db_name
    DB_HOST=your_db_host
    DB_PORT=your_db_port
    ```

4. Run the application:

    ```sh
    go run main.go
    ```

### Running Tests

To run the tests, use the following command:

```sh
go test -v