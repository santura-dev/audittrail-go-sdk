package main

import (
	"fmt"

	"github.com/santura-dev/audittrail-go-sdk/audittrail"
)

// This example demonstrates how to use the AuditTrail Go SDK.
// Before running, ensure you have:
// 1. A running instance of the AuditTrail API (e.g., deployed at https://api.example.com
//    or locally at http://localhost:8080). See the main project at
//    https://github.com/santura-dev/AuditTrail for setup.
// 2. A valid JWT token obtained from the API's token endpoint.
// 3. Go installed (version 1.22 or later).

func main() {
	// Replace these with your own API base URL and JWT token.
	// The baseURL should point to where your AuditTrail API is running.
	baseURL := "https://api.example.com"
	token := "your-jwt-token-here"

	// Initialize the AuditTrail client. This will return an error if the baseURL is invalid.
	client, err := audittrail.NewAuditTrailClient(baseURL, token)
	if err != nil {
		fmt.Printf("Failed to initialize client: %v\n", err)
		return
	}

	// Example 1: Create a new log entry
	response, err := client.CreateLog("login", map[string]interface{}{"ip": "192.168.1.1"})
	if err != nil {
		fmt.Printf("Failed to create log: %v\n", err)
		return
	}
	fmt.Printf("Create Log Response: %+v\n", response)

	// Example 2: List logs with filters
	params := map[string]string{
		"action__contains": "login", // Filter logs containing "login" (case-insensitive regex)
		"page":             "1",     // Page number for pagination
		"page_size":        "20",    // Number of logs per page (max 100)
	}
	logs, err := client.ListLogs(params)
	if err != nil {
		fmt.Printf("Failed to list logs: %v\n", err)
		return
	}
	fmt.Printf("Total Logs: %d\n", logs.Count)
	for i, log := range logs.Results {
		fmt.Printf("Log %d: %+v\n", i+1, log)
	}
}