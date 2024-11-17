# Task Management Application

This repository contains a simplified Task Management Application with both backend and frontend components.  
The goal of this project is to demonstrate REST API design, React-based frontend development, deployment strategies using Docker, and Continuous Integration (CI) setup with GitHub Actions.


## Components
### Backend (Golang)
- REST API for managing tasks.
- Endpoints to create, read, update, and delete tasks.
- In-memory data storage for simplicity.
- Unit-tests for reliability.

### Frontend (React)
- User-friendly interface to interact with tasks.
- Add, edit, view, and delete tasks.
- Minimal and responsive UI.
- Integration with backend API.

### Deployment
- Dockerized backend and frontend.
- `docker-compose` for local setup.
- CI pipeline with GitHub Actions for testing, building, and linting.

---

## Setup Instructions

### Prerequisites
- [Docker](https://www.docker.com/)
- [Node.js](https://nodejs.org/) (optional for local frontend development)
- [Golang](https://go.dev/) (optional for local backend development)

### Local Development and Deployment
1. Clone the repository:
   ```bash
   git clone https://github.com/saarzur123/task-management.git
   cd task-management-app
   ```

2. **Using Docker Compose**:
    - Ensure Docker is running on your machine.
    - Build and start the application:
      ```bash
      docker-compose up --build
      ```
    - The backend API will be available at `http://localhost:8080`.
    - The frontend will be available at `http://localhost:3000`.

3. **Manual Setup (Optional)**:
    - **Backend**:
        - Navigate to the `backend` folder.
        - Run the application:
          ```bash
          go run main.go
          ```
    - **Frontend**:
        - Navigate to the `frontend` folder.
        - Install dependencies:
          ```bash
          npm install
          ```
        - Start the development server:
          ```bash
          npm start
          ```

### API Endpoints
- **POST /tasks**: Create a new task.
- **GET /tasks/{id}**: Retrieve task details by ID.
- **PUT /tasks/{id}**: Update task details.
- **DELETE /tasks/{id}**: Remove a task.

Refer to the API documentation in the `backend` directory for more details.

---

## Testing
- **Backend Tests**:
    - Run tests:
      ```bash
      make backend-test
      ```
- **Frontend Tests**:
    - Run tests:
      ```bash
      make frontend-test
      ```
- **All Tests**:
    - Run tests:
      ```bash
      make test
      ```
---

## CI/CD Pipeline
The project includes a GitHub Actions pipeline with the following features:
1. Automated testing for both backend and frontend on every pull request.
2. Code linting for quality assurance.
3. Docker image builds for backend and frontend.