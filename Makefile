run:
	@echo "Running the application"
	@go mod tidy
	@go run cmd/main.go

build:
	@echo "Building the application"
	go build -o bin/main main.go

clean:
	@echo "Cleaning the application"
	rm -rf bin

watch:
	@if command -v air > /dev/null; then \
            air; \
            echo "Watching...";\
        else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                air; \
                echo "Watching...";\
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi

PHONY: run build clean
	