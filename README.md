# AuditTrail Go SDK

A Go client library for interacting with the AuditTrail API. This SDK provides a simple and reliable way to create and query audit log entries via the AuditTrail API, with built-in support for retries and request tracing.

## Overview

The AuditTrail Go SDK allows developers to easily integrate with the AuditTrail API, which is designed to log and retrieve audit trail data. It handles HTTP requests, JWT authentication, retries, and request IDs, making it suitable for production use.

- **Version**: v0.1.0
- **License**: MIT
- **Repository**: https://github.com/santura-dev/audittrail-go-sdk
- **Main Project**: https://github.com/santura-dev/AuditTrail (for API setup instructions)

## Installation

To use the SDK, you need Go installed (version 1.24.1 or later). Install the SDK using go get:

```bash
go get github.com/santura-dev/audittrail-go-sdk@v0.1.0
```

This will add the SDK as a dependency in your project's go.mod file.

## Prerequisites

Before using the SDK, ensure you have:

1. A running AuditTrail API instance:
   - Deploy your own instance of the AuditTrail API
   - Refer to the main project repository for setup instructions
   - The API must be accessible at a valid URL (e.g., https://api.example.com or http://localhost:8080)

2. A valid JWT token:
   - Obtain a Bearer token from the API's token endpoint (typically `/api/token/` if using Django with rest_framework_simplejwt)
   - The token is required for authentication

## Usage

The SDK provides two main methods: `CreateLog` to add a new log entry and `ListLogs` to retrieve logs with filters. Below is an example demonstrating both.

### Example Code

```go
package main

import (
    "fmt"
    "github.com/santura-dev/audittrail-go-sdk/audittrail"
)

func main() {
    // Replace these with your own API base URL and JWT token
    baseURL := "https://api.example.com" // Update this!
    token := "your-jwt-token-here"      // Update this!

    // Initialize the AuditTrail client
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
        "page":            "1",     // Page number for pagination
        "page_size":       "20",    // Number of logs per page (max 100)
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
```

### Running the Example

1. Clone this repository:
```bash
git clone https://github.com/santura-dev/audittrail-go-sdk.git
cd audittrail-go-sdk
```

2. Update main.go with your baseURL and token
3. Run the example:
```bash
go run main.go
```

Expected output (if successful):
```
Create Log Response: map[message:Log created]
Total Logs: 1
Log 1: {ID:12345 Timestamp:2025-05-18T18:42:00Z Action:login UserID:user123 Details:map[ip:192.168.1.1] Signature:abc123...}
```

## API Methods

### NewAuditTrailClient(baseURL, token string) (*AuditTrailClient, error)
- **Purpose**: Initializes a new client instance
- **Parameters**:
  - `baseURL`: The URL where the AuditTrail API is running (e.g., https://api.example.com)
  - `token`: A valid JWT token for authentication
- **Returns**: A *AuditTrailClient or an error if the baseURL is invalid or missing

### CreateLog(action string, details map[string]interface{}) (map[string]string, error)
- **Purpose**: Creates a new log entry
- **Parameters**:
  - `action`: The action to log (e.g., "login")
  - `details`: Optional additional details as a map (e.g., map[string]interface{}{"ip": "192.168.1.1"})
- **Returns**: A response map (e.g., map[string]string{"message": "Log created"}) or an error

### ListLogs(params map[string]string) (*LogListResponse, error)
- **Purpose**: Retrieves a paginated list of logs with filters
- **Parameters**:
  - `params`: A map of query parameters, e.g.:
    - `action__contains`: Filter by partial action name (regex)
    - `page`: Page number
    - `page_size`: Number of logs per page (max 100)
    - Other supported filters: user_id, action, action__in, action__nin, start_time, end_time
- **Returns**: A *LogListResponse (with Count, Next, Previous, and Results) or an error

## Features

- **Retries**: Automatically retries failed requests up to 3 times with exponential backoff (1s to 10s)
- **Request Tracing**: Adds a unique X-Request-ID header to each request for debugging
- **Error Handling**: Returns detailed errors for invalid requests or server issues


## License

This project is licensed under the MIT License - see the LICENSE file for details.


## Acknowledgements

- Built as part of the AuditTrail project by santura-dev
- Uses resty for HTTP requests and google/uuid for request IDs
