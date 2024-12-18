BACKEND_DIR=backend
FRONTEND_DIR=frontend

.PHONY: test lint backend-test frontend-test backend-lint frontend-lint docker-backend docker-frontend

# Run tests
test: backend-test frontend-test

backend-test:
	go clean -testcache
	@echo "Running backend tests..."
	cd $(BACKEND_DIR) && go test ./...
	@echo "Done"

frontend-test:
	@echo "Running frontend tests..."
	cd $(FRONTEND_DIR) && npm install && npm test -- --watchAll=false
	@echo "Done"

# Lint
lint: backend-lint frontend-lint

backend-lint:
	@echo "Linting backend code..."
	 cd $(BACKEND_DIR) && curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b /usr/local/bin v1.61.0 && ./bin/golangci-lint run --config golangci.yml --skip-dirs "/usr/local|/opt/homebrew/Cellar/go"
	@echo "Done"

frontend-lint:
	@echo "Linting frontend code..."
	cd $(FRONTEND_DIR) && npx eslint src --ext .js,.jsx,.ts,.tsx
	@echo "Done"

# Build Docker images
docker-backend:
	@echo "Building backend Docker image..."
	docker build -t backend-service:latest $(BACKEND_DIR)

docker-frontend:
	@echo "Building frontend Docker image..."
	docker build -t frontend-service:latest $(FRONTEND_DIR)

