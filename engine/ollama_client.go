package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// OllamaClient handles communication with the Ollama API
type OllamaClient struct {
	baseURL string
	client  *http.Client
}

// GenerateRequest represents a request to the Ollama generate API
type GenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

// GenerateResponse represents a response from the Ollama generate API
type GenerateResponse struct {
	Model     string    `json:"model"`
	CreatedAt time.Time `json:"created_at"`
	Response  string    `json:"response"`
	Done      bool      `json:"done"`
}

// HealthResponse represents a response from the Ollama health check
type HealthResponse struct {
	Status string `json:"status"`
}

// NewOllamaClient creates a new Ollama API client
func NewOllamaClient(serverAddr string) *OllamaClient {
	return &OllamaClient{
		baseURL: fmt.Sprintf("http://%s", serverAddr),
		client: &http.Client{
			Timeout: 300 * time.Second, // 5 minute timeout for LLM responses
		},
	}
}

// HealthCheck verifies the Ollama server is accessible
func (c *OllamaClient) HealthCheck() error {
	resp, err := c.client.Get(c.baseURL + "/api/tags")
	if err != nil {
		return fmt.Errorf("failed to connect to Ollama server at %s: %v", c.baseURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Ollama server returned status %d", resp.StatusCode)
	}

	return nil
}

// Generate sends a prompt to the specified model and returns the response
func (c *OllamaClient) Generate(model, prompt string) (string, error) {
	req := GenerateRequest{
		Model:  model,
		Prompt: prompt,
		Stream: false, // Use non-streaming for simplicity
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	resp, err := c.client.Post(
		c.baseURL+"/api/generate",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	var response GenerateResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	return response.Response, nil
}

// GenerateStream sends a prompt and returns a channel for streaming responses
func (c *OllamaClient) GenerateStream(model, prompt string) (<-chan string, <-chan error) {
	responses := make(chan string)
	errors := make(chan error, 1)

	go func() {
		defer close(responses)
		defer close(errors)

		req := GenerateRequest{
			Model:  model,
			Prompt: prompt,
			Stream: true,
		}

		jsonData, err := json.Marshal(req)
		if err != nil {
			errors <- fmt.Errorf("failed to marshal request: %v", err)
			return
		}

		resp, err := c.client.Post(
			c.baseURL+"/api/generate",
			"application/json",
			bytes.NewBuffer(jsonData),
		)
		if err != nil {
			errors <- fmt.Errorf("failed to send request: %v", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			errors <- fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
			return
		}

		decoder := json.NewDecoder(resp.Body)
		for {
			var response GenerateResponse
			if err := decoder.Decode(&response); err != nil {
				if err == io.EOF {
					break
				}
				errors <- fmt.Errorf("failed to decode response: %v", err)
				return
			}

			responses <- response.Response

			if response.Done {
				break
			}
		}
	}()

	return responses, errors
}

// ListModels returns the list of available models
func (c *OllamaClient) ListModels() ([]string, error) {
	resp, err := c.client.Get(c.baseURL + "/api/tags")
	if err != nil {
		return nil, fmt.Errorf("failed to get models: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var result struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	var models []string
	for _, model := range result.Models {
		models = append(models, model.Name)
	}

	return models, nil
}