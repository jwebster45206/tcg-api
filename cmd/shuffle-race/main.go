package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	// API base URL - hardcoded assumption for Docker environment
	APIBaseURL = "http://localhost:8080"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "help", "-h", "--help":
		printUsage()
		os.Exit(0)

	case "health":
		checkAPIHealth()
		os.Exit(0)

	case "race":
		decks := 3
		if len(os.Args) > 2 {
			if d, err := strconv.Atoi(os.Args[2]); err == nil {
				decks = d
			} else {
				fmt.Printf("Invalid number of decks: %s\n", os.Args[2])
				os.Exit(1)
			}
		}
		fmt.Println("Starting shuffle..., decks:", decks)
		if err := race(decks); err != nil {
			fmt.Printf("Error occurred during shuffle: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)

	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: go run *.go <command>")
	fmt.Println("Commands:")
	fmt.Println("  help     Show this help message")
	fmt.Println("  version  Show version information")
	fmt.Println("  health   Check API health status")
	fmt.Println("  race     Run a shuffle race experiment")
}

// HealthResponse represents the expected response from the health endpoint
type HealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

func newHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second,
	}
}

// checkAPIHealth connects to the API health endpoint and reports status
func checkAPIHealth() {
	// Hardcoded assumption: running in Docker, API is accessible at tcg-api:8080
	apiURL := APIBaseURL + "/health"

	fmt.Println("Checking API health...")
	fmt.Printf("Connecting to: %s\n", apiURL)

	client := newHTTPClient()

	// Make the request
	resp, err := client.Get(apiURL)
	if err != nil {
		fmt.Printf("❌ Failed to connect to API: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = resp.Body.Close() }()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("❌ API returned non-OK status: %d\n", resp.StatusCode)
		os.Exit(1)
	}

	// Parse JSON response
	var health HealthResponse
	if err := json.NewDecoder(resp.Body).Decode(&health); err != nil {
		fmt.Printf("❌ Failed to parse health response: %v\n", err)
		os.Exit(1)
	}

	// Check health status
	if health.Status == "ok" || health.Status == "healthy" {
		fmt.Println("✅ API is healthy and responding!")
		fmt.Printf("Status: %s\n", health.Status)
		if health.Message != "" {
			fmt.Printf("Message: %s\n", health.Message)
		}
	} else {
		fmt.Printf("⚠️  API responded but status is not healthy: %s\n", health.Status)
		if health.Message != "" {
			fmt.Printf("Message: %s\n", health.Message)
		}
		os.Exit(1)
	}
}
