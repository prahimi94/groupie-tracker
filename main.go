package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	// Print the current working directory for debugging
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Current working directory:", cwd)

	// Start the backend service (running main.go in backend/api)
	fmt.Println("Starting backend...")
	cmdBackend := exec.Command("go", "run", "backend/api/main.go")
	cmdBackend.Stdout = os.Stdout
	cmdBackend.Stderr = os.Stderr

	// Start the second backend service (running go-routine/main.go)
	fmt.Println("Starting backend with go-routine...")
	cmdBackendRoutine := exec.Command("go", "run", "backend/api/go-routine/main.go")
	cmdBackendRoutine.Stdout = os.Stdout
	cmdBackendRoutine.Stderr = os.Stderr

	// Run both commands concurrently
	errChan := make(chan error, 2)

	go func() {
		errChan <- cmdBackend.Run()
	}()

	go func() {
		errChan <- cmdBackendRoutine.Run()
	}()

	// Wait for both commands to finish or log an error if any fails
	for i := 0; i < 2; i++ {
		if err := <-errChan; err != nil {
			log.Fatalf("Error running a backend command: %v", err)
		}
	}
}
