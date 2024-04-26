# Todo REST API(V2)

This is a simple REST API built using Go and Mysql. The API provides endpoints for user authentication (registration and login) and CRUD operations for managing todo items. The goal of this project was to stick to the Go standard library as much as possible, avoiding the use of external libraries or frameworks where possible.

## Technologies Used

- **Go** (version 1.22)
- **MySQL** (database)
- **Docker** (for hosting the MySQL database)
- **JWT** (for authentication)

### External Packages

- [golang-jwt/jwt](github.com/golang-jwt/jwt/v5) (for JWT implementation)
- [go-sql-driver/mysql](github.com/go-sql-driver/mysql) (for MySQL database driver)
- [joho/godotenv](github.com/joho/godotenv) (for loading environment variables from `.env` file)

### Folder Structure
```go
go-todo/
├── handlers/
│   ├── auth.go
│   ├── handlers.go
│   └── todo.go
├── models/
│   ├── response.go
│   ├── todo.go
│   └── user.go
├── services/
│   ├── db.go
│   ├── todo.go
│   └── user.go
├── utils/
│   └── utils.go
├── .env
├── .gitignore
├── go.mod
├── go.sum
└── main.go
```
- `handlers/`: Contains handler functions for different routes.
  - `auth.go`: Handlers for user authentication (login, register) and authentication middleware.
  - `handlers.go`: Serves the handlers and applies middleware.
  - `todo.go`: Handlers for todo CRUD operations.
- `models/`: Contains data models for the application.
  - `response.go`: Response structs for API responses.
  - `todo.go`: Todo model struct.
  - `user.go`: User model struct.
- `services/`: Contains service functions for database operations.
  - `db.go`: Initializes the database connection.
  - `todo.go`: CRUD functions for todo items.
  - `user.go`: CRUD functions for user authentication.
- `utils/`: Contains utility functions.
  - `utils.go`: Utility functions, including JWT token creation and verification.
- `main.go`: Entry point of the application, imports and serves the handlers.
- `.env`: Environment variable file (not committed to version control).
- `.gitignore`: Git ignore file.
- `go.mod` and `go.sum`: Go module files.

### Routes

#### Authentication Routes

- `POST /api/register`: Register a new user.
- `POST /api/login`: Authenticate a user and obtain a JWT token.

#### Todo Routes

- `POST /api/todo`: Create a new todo item (requires authentication).
- `GET /api/todo/{id}`: Retrieve a specific todo item (requires authentication).
- `GET /api/todos`: Retrieve all todo items (requires authentication).
- `DELETE /api/todo/{id}`: Delete a specific todo item (requires authentication).
- `PUT /api/todo/{id}`: Update a specific todo item (requires authentication).

### Authentication

The API uses JWT (JSON Web Tokens) for authentication. After a successful login, the API returns a JWT token in a cookie, which should be included in subsequent requests for authenticated routes.

### Database

The application uses a MySQL database for storing user and todo data. The database connection is initialized in `services/db.go`. The `User` and `Todo` tables were created manually by running the following SQL queries using a SQL client (in this case, TablePlus):

```sql
CREATE TABLE User (
    id int NOT NULL AUTO_INCREMENT,
    name varchar(255) NOT NULL,
    email varchar(255),
    role varchar(255),
    password varchar(255),
    PRIMARY KEY (id)
);
```
```sql
CREATE TABLE Todo (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    user_id INT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES User(id)
);
```

### Middleware

The API uses a custom middleware function `ApplyMiddleware` to apply middleware functions (e.g., `AuthMiddleware`) to specific routes. The `AuthMiddleware` function retrieves the JWT token from the cookie, verifies it, and forwards the request if the token is valid for authenticated routes.

### Error Handling

The API handles errors and returns appropriate HTTP status codes and error messages in the response.

### Environment Variables

The application uses environment variables for configuration. The required environment variables are stored in the `.env` file (which is not committed to version control). The environment variables include:

- `DB_HOST`: MySQL database host
- `DB_PORT`: MySQL database port
- `DB_USERNAME`: MySQL database user
- `DB_PASSWORD`: MySQL database password
- `DB_NAME`: MySQL database name
- `JWT_SECRET`: Secret key for JWT signing and verification

### Running the Application

To run the Go Todo REST API, follow these steps:

1. Set up a Docker container for MySQL.
2. Run the above mentioned SQL queries manually to create the `User` and `Todo` tables.
3. Update the environment variables in the `.env` file with the correct MySQL database connection details.
4. Install the required Go packages: `go get ./...`
5. Run the application: `go run main.go`

The API will be served on `http://localhost:8080`.
