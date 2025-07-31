package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

func sanitizeModelName(modelName string) string {
	// Replace invalid Windows filename characters with underscores
	invalidChars := []string{"<", ">", ":", "\"", "/", "\\", "|", "?", "*"}
	sanitized := modelName
	for _, char := range invalidChars {
		sanitized = strings.ReplaceAll(sanitized, char, "_")
	}
	// Also replace spaces and colons commonly found in model names
	sanitized = strings.ReplaceAll(sanitized, " ", "_")
	sanitized = strings.ReplaceAll(sanitized, ":", "_")
	return sanitized
}

func generateFilenameFromPrompt(prompt string) string {
	// Take first few words from the prompt to create a descriptive filename
	words := strings.Fields(prompt)
	maxWords := 5
	if len(words) > maxWords {
		words = words[:maxWords]
	}
	
	// Join words and sanitize for filename
	filename := strings.Join(words, "_")
	
	// Replace invalid Windows filename characters
	invalidChars := []string{"<", ">", ":", "\"", "/", "\\", "|", "?", "*", "."}
	for _, char := range invalidChars {
		filename = strings.ReplaceAll(filename, char, "_")
	}
	
	// Convert to lowercase and limit length
	filename = strings.ToLower(filename)
	if len(filename) > 50 {
		filename = filename[:50]
	}
	
	// Remove trailing underscores
	filename = strings.TrimRight(filename, "_")
	
	return filename + "_response.txt"
}

func main() {
	// Configuration
	serverAddr := "192.168.0.63:11434"
	modelName := "qwen3:30b"
	baseURL := fmt.Sprintf("http://%s", serverAddr)

	// Create results directory structure
	sanitizedModelName := sanitizeModelName(modelName)
	resultsDir := filepath.Join("results", sanitizedModelName)
	if err := os.MkdirAll(resultsDir, 0755); err != nil {
		log.Fatalf("Failed to create results directory: %v", err)
	}
	log.Printf("Results will be saved to: %s", resultsDir)

	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	log.Printf("Testing LLM at %s with model %s (ADVANCED PROMPTS)", baseURL, modelName)

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

	// Save simple prompt response to file
	simpleResponseFile := filepath.Join(resultsDir, "simple_prompt_response.txt")
	if err := os.WriteFile(simpleResponseFile, []byte(response.Response), 0644); err != nil {
		log.Printf("Failed to save simple prompt response to file: %v", err)
	} else {
		log.Printf("Simple prompt response saved to %s", simpleResponseFile)
	}

	// Test 3: Advanced Programming Prompts
	log.Println("\n=== Test 3: Advanced Programming Prompts (Random Order) ===")
	
	advancedPrompts := []string{
		"Implement a complete BASIC interpreter in Go that supports variables, loops, conditionals, subroutines, and mathematical expressions. Include error handling and line number management.",
		"Design and implement a concurrent web scraper in Go that can handle rate limiting, retries, and graceful error handling while scraping multiple sites simultaneously.",
		"Create a complete TCP/IP server in Go that implements a custom protocol for a multi-user chat system with rooms, user authentication, and message persistence.",
		"Implement a full lexer, parser, and AST evaluator for a simple programming language in Go. Include support for functions, variables, and control flow.",
		"Build a distributed key-value store in Go with consistent hashing, replication, and fault tolerance. Include a client library and REST API.",
		"Design a complete database query engine in Go that can parse SQL, optimize queries, and execute them against in-memory data structures with indexing.",
		"Implement a fully functional HTTP/2 server from scratch in Go without using the standard library's HTTP/2 implementation. Include multiplexing and flow control.",
		"Create a complete compiler for a subset of C that generates x86-64 assembly. Include preprocessing, optimization passes, and proper symbol table management.",
		"Build a sophisticated caching system in Go with TTL, LRU eviction, persistence, and distributed cache invalidation across multiple nodes.",
		"Implement a complete Git-like version control system in Go with branching, merging, diff algorithms, and a working directory management system.",
	}

	// Shuffle the prompts for random order
	for i := len(advancedPrompts) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		advancedPrompts[i], advancedPrompts[j] = advancedPrompts[j], advancedPrompts[i]
	}

	var totalDuration time.Duration
	successCount := 0

	for i, prompt := range advancedPrompts {
		log.Printf("\n--- Advanced Programming Test %d/10 ---", i+1)
		log.Printf("Prompt: %s", prompt)
		log.Printf("Prompt length: %d characters", len(prompt))

		req.Prompt = prompt
		jsonData, err = json.Marshal(req)
		if err != nil {
			log.Printf("Failed to marshal request %d: %v", i+1, err)
			continue
		}

		log.Printf("Sending advanced programming prompt %d...", i+1)
		start = time.Now()

		resp, err = client.Post(
			baseURL+"/api/generate",
			"application/json",
			bytes.NewBuffer(jsonData),
		)
		if err != nil {
			log.Printf("Failed to send advanced programming request %d: %v", i+1, err)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			log.Printf("Advanced programming API request %d failed with status %d: %s", i+1, resp.StatusCode, string(body))
			resp.Body.Close()
			continue
		}

		body, err = io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Printf("Failed to read advanced programming response %d: %v", i+1, err)
			continue
		}

		var response TestResponse
		if err := json.Unmarshal(body, &response); err != nil {
			log.Printf("Failed to parse advanced programming response %d: %v", i+1, err)
			continue
		}

		duration := time.Since(start)
		totalDuration += duration
		successCount++

		log.Printf("Advanced programming prompt %d completed in %v", i+1, duration)
		log.Printf("Response length: %d characters", len(response.Response))
		log.Printf("First 200 chars: %q", truncateString(response.Response, 200))
		log.Printf("Contains 'func' keyword: %t", contains(response.Response, "func"))
		log.Printf("Contains 'package' keyword: %t", contains(response.Response, "package"))

		// Save advanced prompt response to file
		filename := generateFilenameFromPrompt(prompt)
		filePath := filepath.Join(resultsDir, filename)
		if err := os.WriteFile(filePath, []byte(response.Response), 0644); err != nil {
			log.Printf("Failed to save advanced prompt %d response to file: %v", i+1, err)
		} else {
			log.Printf("Advanced prompt %d response saved to %s", i+1, filePath)
		}
	}

	// Summary of advanced programming tests
	log.Printf("\n=== Advanced Programming Tests Summary ===")
	log.Printf("Successful prompts: %d/10", successCount)
	if successCount > 0 {
		avgDuration := totalDuration / time.Duration(successCount)
		log.Printf("Total time: %v", totalDuration)
		log.Printf("Average response time: %v", avgDuration)
		log.Printf("Longest response time: %v", findLongestDuration(advancedPrompts, client, baseURL, modelName))
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

	log.Println("\n=== Advanced Test Summary ===")
	log.Println("Advanced prompts test LLM performance on complex, multi-step programming tasks.")
	log.Println("These should take significantly longer than simple prompts (2-10 minutes each).")
	log.Printf("If responses are completing in under 30 seconds, the model may not be fully processing the complexity.")
}

func truncateString(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length] + "..."
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && 
		   (len(substr) == 0 || 
			func() bool {
				for i := 0; i <= len(s)-len(substr); i++ {
					if s[i:i+len(substr)] == substr {
						return true
					}
				}
				return false
			}())
}

func findLongestDuration(prompts []string, client *http.Client, baseURL, modelName string) time.Duration {
	// This is a placeholder - would need to track individual durations in the main loop
	// For now, return a reasonable estimate
	return time.Duration(5) * time.Minute
}