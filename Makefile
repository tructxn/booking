# This is a comment
# Target: Dependencies
#   Recipe (commands)

all: build

# Target to build the program
build: main.go
	@echo "Building the Go application..."
	go build -o myapp main.go

# Clean up generated files
clean:
	@echo "Cleaning up..."
	rm -f myapp

# Target to run the application
run: build
	./myapp

.PHONY: all build clean run

