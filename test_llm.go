package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type TestRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type TestResponse struct {
	Model     string    `json:"model"`
	CreatedAt time.Time `json:"created_at"`
	Response  string    `json:"response"`
	Done      bool      `json:"done"`
}

func main() {
	// Configuration
	serverAddr := "192.168.0.63:11434"
	modelName := "qwen3:30b"
	baseURL := fmt.Sprintf("http://%s", serverAddr)

	log.Printf("Testing LLM at %s with model %s", baseURL, modelName)

	// Create HTTP client with no timeout to see how long it actually takes
	client := &http.Client{
		Timeout: 0, // No timeout
	}

	// Test 1: Health check
	log.Println("=== Test 1: Health Check ===")
	start := time.Now()
	resp, err := client.Get(baseURL + "/api/tags")
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	resp.Body.Close()
	log.Printf("Health check completed in %v (status: %d)", time.Since(start), resp.StatusCode)

	// Test 2: Simple prompt
	log.Println("\n=== Test 2: Simple Prompt ===")
	simplePrompt := "Hello, what is 2+2?"
	
	req := TestRequest{
		Model:  modelName,
		Prompt: simplePrompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		log.Fatalf("Failed to marshal request: %v", err)
	}

	log.Printf("Sending simple prompt: %q", simplePrompt)
	start = time.Now()

	resp, err = client.Post(
		baseURL+"/api/generate",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Fatalf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}

	var response TestResponse
	if err := json.Unmarshal(body, &response); err != nil {
		log.Fatalf("Failed to parse response: %v", err)
	}

	duration := time.Since(start)
	log.Printf("Simple prompt completed in %v", duration)
	log.Printf("Response length: %d characters", len(response.Response))
	log.Printf("Response: %q", response.Response)

	// Test 3: Programming prompt (similar to what the engine sends)
	log.Println("\n=== Test 3: Programming Prompt ===")
	codingPrompt := `You are an expert software developer. Your task is to implement a simple BASIC interpreter in Go. 

Requirements:
1. Support line-numbered BASIC syntax
2. Implement PRINT and LET statements
3. Include basic error handling

Please provide a minimal Go implementation.`

	req.Prompt = codingPrompt
	jsonData, err = json.Marshal(req)
	if err != nil {
		log.Fatalf("Failed to marshal request: %v", err)
	}

	log.Printf("Sending programming prompt (%d chars)", len(codingPrompt))
	start = time.Now()

	resp, err = client.Post(
		baseURL+"/api/generate",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		log.Fatalf("Failed to send programming request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Fatalf("Programming API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read programming response: %v", err)
	}

	if err := json.Unmarshal(body, &response); err != nil {
		log.Fatalf("Failed to parse programming response: %v", err)
	}

	duration = time.Since(start)
	log.Printf("Programming prompt completed in %v", duration)
	log.Printf("Response length: %d characters", len(response.Response))
	log.Printf("First 200 chars: %q", truncateString(response.Response, 200))

	// Test 4: Model info
	log.Println("\n=== Test 4: Model Information ===")
	modelReq := map[string]string{"name": modelName}
	jsonData, _ = json.Marshal(modelReq)

	start = time.Now()
	resp, err = client.Post(
		baseURL+"/api/show",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		log.Printf("Failed to get model info: %v", err)
	} else {
		defer resp.Body.Close()
		log.Printf("Model info request completed in %v (status: %d)", time.Since(start), resp.StatusCode)
		if resp.StatusCode == http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			var modelInfo map[string]interface{}
			if json.Unmarshal(body, &modelInfo) == nil {
				if license, ok := modelInfo["license"].(string); ok {
					log.Printf("Model license: %s", license)
				}
				if size, ok := modelInfo["size"].(float64); ok {
					log.Printf("Model size: %.2f GB", size/1e9)
				}
			}
		}
	}

	log.Println("\n=== Test Summary ===")
	log.Println("If you see this message, the LLM server is responding normally.")
	log.Println("Compare the response times above to identify any performance issues.")
}

func truncateString(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length] + "..."
}