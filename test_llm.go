package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
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

	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

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

	// Test 3: Multiple Programming Prompts
	log.Println("\n=== Test 3: Programming Prompts (Random Order) ===")
	
	programmingPrompts := []string{
		"Write a simple Go function to calculate factorial of a number.",
		"Create a Go program that reverses a string without using built-in functions.",
		"Implement a basic stack data structure in Go with push and pop operations.",
		"Write a Go function to check if a number is prime.",
		"Create a Go program that finds the largest element in an array.",
		"Implement a simple binary search function in Go.",
		"Write a Go function to count vowels in a string.",
		"Create a Go program that sorts an array using bubble sort.",
		"Implement a basic queue data structure in Go.",
		"Write a Go function to calculate the nth Fibonacci number.",
	}

	// Shuffle the prompts for random order
	for i := len(programmingPrompts) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		programmingPrompts[i], programmingPrompts[j] = programmingPrompts[j], programmingPrompts[i]
	}

	var totalDuration time.Duration
	successCount := 0

	for i, prompt := range programmingPrompts {
		log.Printf("\n--- Programming Test %d/10 ---", i+1)
		log.Printf("Prompt: %s", prompt)
		log.Printf("Prompt length: %d characters", len(prompt))

		req.Prompt = prompt
		jsonData, err = json.Marshal(req)
		if err != nil {
			log.Printf("Failed to marshal request %d: %v", i+1, err)
			continue
		}

		log.Printf("Sending programming prompt %d...", i+1)
		start = time.Now()

		resp, err = client.Post(
			baseURL+"/api/generate",
			"application/json",
			bytes.NewBuffer(jsonData),
		)
		if err != nil {
			log.Printf("Failed to send programming request %d: %v", i+1, err)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			log.Printf("Programming API request %d failed with status %d: %s", i+1, resp.StatusCode, string(body))
			resp.Body.Close()
			continue
		}

		body, err = io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Printf("Failed to read programming response %d: %v", i+1, err)
			continue
		}

		var response TestResponse
		if err := json.Unmarshal(body, &response); err != nil {
			log.Printf("Failed to parse programming response %d: %v", i+1, err)
			continue
		}

		duration := time.Since(start)
		totalDuration += duration
		successCount++

		log.Printf("Programming prompt %d completed in %v", i+1, duration)
		log.Printf("Response length: %d characters", len(response.Response))
		log.Printf("First 150 chars: %q", truncateString(response.Response, 150))
	}

	// Summary of programming tests
	log.Printf("\n=== Programming Tests Summary ===")
	log.Printf("Successful prompts: %d/10", successCount)
	if successCount > 0 {
		avgDuration := totalDuration / time.Duration(successCount)
		log.Printf("Total time: %v", totalDuration)
		log.Printf("Average response time: %v", avgDuration)
	}

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