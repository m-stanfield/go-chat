# Simple Makefile for a Go project

# Build the application
all: build test


backend:
	@cd backend && go run cmd/api/main.go 


front:
	@npm install --prefix ./frontend
	@npm run dev --prefix ./frontend

# Test the application
test:
	@echo "Testing..."
	@cd backend && go test ./... -v


# Live Reload
watch:
	@if command -v air > /dev/null; then \
            cd backend && air; \
            echo "Watching...";\
        else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                cd backend && air; \
                echo "Watching...";\
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi

.PHONY: all build run test clean watch front backend

