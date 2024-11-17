# Design Overview
This directory implements a REST API for managing tasks, structured using a combination of components to handle different layers of the application.  
These components include:  
- **Handler**: Responsible for interacting with the Service layer and defining the API endpoints.
- **Service**: Handles business logic and data access, interacting with the database.  

The chosen database is a lightweight, in-memory **SQLite** database, which is initialized and closed within the main function.

## Task Model
Each task is represented by the following model:

```go
    type Task struct {
        ID          string    `json:"id"`
        Title       string    `json:"title"`
        Description string    `json:"description"`
        Status      string    `json:"status"`
        CreatedAt   time.Time `json:"created_at"`
    }
```

## Components
### Handler
The `Handler` defines the `TaskHandler` struct, which is responsible for processing incoming HTTP requests and invoking the corresponding methods in the Service layer.

```go
type TaskHandler struct {
    DB service.TaskRepository
}
```
The `TaskHandler` struct includes a reference to the `service.TaskRepository` and implements the API endpoints for CRUD operations. Each handler method accepts `http.ResponseWriter` and `http.Request` as parameters, and returns an appropriate response status based on the outcome of the request.  

The available API endpoints are:

- **CREATE/GET/OPTIONS**: http://localhost:8080/tasks
- **GET/UPDATE/DELETE/OPTIONS**: http://localhost:8080/tasks/{id}

### Service
The `Service` layer defines the `TaskRepository` interface, which corresponds to the CRUD operations required by the API.  
The `TaskManager` struct is defined within the service package, responsible for executing the necessary database queries:

```go
type TaskManager struct {
    DB *sql.DB
}
```
The `TaskManager` struct interacts directly with the SQLite database using its `DB` property to perform CRUD operations as part of implementing the `TaskRepository` interface.

# Setup
### Prerequisites
- [Golang](https://go.dev/) (optional for local backend development)

### Run locally
- Navigate to the `backend` folder.
  - Run the application:
  ```bash
  go run main.go
  ```


## Pros and Cons of the Given Implementation
### Pros
1. **Separation of Concerns (SoC)**:  
The architecture separates concerns effectively by using distinct layers: Handler for HTTP request handling and Service for business logic. This separation makes the codebase easier to maintain, test, and extend.
This helps in making the code more modular and promotes better organization.
2. **Lightweight and In-Memory Database**:  
The use of an in-memory SQLite database reduces the overhead of setting up and maintaining a persistent database, making the application more lightweight and faster to start.
It stays alive as long as the app is running.
3.**Easy to Test**:
The decoupling between the handler and service layer makes it easier to write unit tests for both the HTTP layer and the business logic layer separately.
Mocking the TaskRepository interface during testing allows the API logic to be tested independently from the database.

### Cons
1. **In-Memory Database Limitation**:
Using an in-memory SQLite database is ideal for testing or development, but it has significant limitations for production environments. It does not persist data between application restarts, meaning any data added during runtime will be lost when the application stops.
For production use, a more robust, persistent database (e.g., PostgreSQL, MySQL) should be used.