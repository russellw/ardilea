package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// Config holds the engine configuration
type Config struct {
	OllamaServer string `json:"ollama_server"`
	ModelName    string `json:"model_name"`
	WorkspaceDir string `json:"workspace_dir"`
}

// Engine represents the LLM agent engine
type Engine struct {
	config *Config
	client *OllamaClient
}

// NewEngine creates a new engine instance
func NewEngine() (*Engine, error) {
	config, err := loadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}

	client := NewOllamaClient(config.OllamaServer)

	return &Engine{
		config: config,
		client: client,
	}, nil
}

// loadConfig reads configuration from config.json with defaults
func loadConfig() (*Config, error) {
	config := &Config{
		OllamaServer: "192.168.0.63:11434",
		ModelName:    "qwen2.5:32b",
		WorkspaceDir: "/workspace",
	}

	configPath := "config.json"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Printf("Config file %s not found, using defaults", configPath)
		return config, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	if err := json.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	log.Printf("Loaded config: Ollama=%s, Model=%s, Workspace=%s", 
		config.OllamaServer, config.ModelName, config.WorkspaceDir)

	return config, nil
}

// Run starts the engine and begins the BASIC interpreter development session
func (e *Engine) Run() error {
	log.Println("Starting LLM Agent Engine...")
	
	// Ensure workspace directory exists
	if err := os.MkdirAll(e.config.WorkspaceDir, 0755); err != nil {
		return fmt.Errorf("failed to create workspace directory: %v", err)
	}

	// Check if we can connect to Ollama
	if err := e.client.HealthCheck(); err != nil {
		return fmt.Errorf("failed to connect to Ollama server: %v", err)
	}

	// Start the development session
	return e.startDevelopmentSession()
}

// startDevelopmentSession begins the interactive development process
func (e *Engine) startDevelopmentSession() error {
	log.Println("Starting BASIC interpreter development session...")

	// Check if BASIC interpreter already exists
	basicPath := filepath.Join(e.config.WorkspaceDir, "basic")

	if _, err := os.Stat(basicPath); err == nil {
		log.Println("BASIC interpreter already exists, analyzing current state...")
		return e.analyzeExistingCode()
	}

	log.Println("No BASIC interpreter found, starting fresh development...")
	return e.startFreshDevelopment()
}

// analyzeExistingCode examines the current workspace and suggests improvements
func (e *Engine) analyzeExistingCode() error {
	// Read the current workspace state
	workspaceFiles, err := e.scanWorkspace()
	if err != nil {
		return fmt.Errorf("failed to scan workspace: %v", err)
	}

	prompt := fmt.Sprintf(`You are an expert software developer assistant. I have a workspace with a BASIC interpreter implementation. Please analyze the current state and suggest next steps for improvement.

Current workspace files:
%s

The goal is to have a complete, well-tested BASIC interpreter. Please:
1. Analyze the current implementation
2. Identify any gaps or areas for improvement  
3. Suggest specific next steps
4. Prioritize the most important improvements

Please be specific and actionable in your suggestions.`, workspaceFiles)

	response, err := e.client.Generate(e.config.ModelName, prompt)
	if err != nil {
		return fmt.Errorf("failed to get LLM response: %v", err)
	}

	log.Println("=== LLM Analysis ===")
	fmt.Println(response)
	log.Println("=== End Analysis ===")

	return nil
}

// startFreshDevelopment begins developing a BASIC interpreter from scratch
func (e *Engine) startFreshDevelopment() error {
	prompt := `You are an expert software developer. Your task is to implement a BASIC interpreter in Go with the following requirements:

1. Support line-numbered BASIC syntax (classic style)
2. Implement core statements: PRINT, LET, GOTO, IF-THEN, FOR-NEXT, REM, END
3. Support variables (both numeric and string)
4. Include proper error handling
5. Accept filename as command line argument

The interpreter should be compatible with test files that exist in tests/basic/ directory.

Please provide a complete Go implementation of the BASIC interpreter. Focus on correctness and clarity.`

	response, err := e.client.Generate(e.config.ModelName, prompt)
	if err != nil {
		return fmt.Errorf("failed to get LLM response: %v", err)
	}

	log.Println("=== LLM Generated Code ===")
	fmt.Println(response)
	log.Println("=== End Generated Code ===")

	// TODO: Parse the response and extract code to write to files
	// TODO: Run tests to verify the generated code
	// TODO: Iterate on improvements

	return nil
}

// scanWorkspace reads the current workspace structure
func (e *Engine) scanWorkspace() (string, error) {
	var result string

	err := filepath.Walk(e.config.WorkspaceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip hidden files and directories
		if filepath.Base(path)[0] == '.' {
			return nil
		}

		relPath, _ := filepath.Rel(e.config.WorkspaceDir, path)
		if info.IsDir() {
			result += fmt.Sprintf("📁 %s/\n", relPath)
		} else {
			size := info.Size()
			result += fmt.Sprintf("📄 %s (%d bytes)\n", relPath, size)
		}

		return nil
	})

	return result, err
}

func main() {
	engine, err := NewEngine()
	if err != nil {
		log.Fatalf("Failed to create engine: %v", err)
	}

	if err := engine.Run(); err != nil {
		log.Fatalf("Engine failed: %v", err)
	}
}