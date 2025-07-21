package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Config holds the engine configuration
type Config struct {
	OllamaServer string `json:"ollama_server"`
	ModelName    string `json:"model_name"`
	WorkspaceDir string `json:"workspace_dir"`
}

// FileInfo represents information about a file
type FileInfo struct {
	Path    string    `json:"path"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"mod_time"`
	Hash    string    `json:"hash"`
	IsDir   bool      `json:"is_dir"`
}

// WorkspaceSnapshot represents the state of the workspace at a point in time
type WorkspaceSnapshot struct {
	Timestamp time.Time            `json:"timestamp"`
	Files     map[string]FileInfo  `json:"files"`
}

// WorkspaceReport compares before and after snapshots
type WorkspaceReport struct {
	Before   WorkspaceSnapshot `json:"before"`
	After    WorkspaceSnapshot `json:"after"`
	Added    []string          `json:"added"`
	Removed  []string          `json:"removed"`
	Modified []string          `json:"modified"`
	Summary  string            `json:"summary"`
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
		ModelName:    "qwen3:30b",
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

	// Take a snapshot before starting
	log.Println("Creating workspace snapshot before engine run...")
	beforeSnapshot, err := e.takeWorkspaceSnapshot()
	if err != nil {
		return fmt.Errorf("failed to create before snapshot: %v", err)
	}

	// Start the development session
	err = e.startDevelopmentSession()

	// Take a snapshot after completion (regardless of success/failure)
	log.Println("Creating workspace snapshot after engine run...")
	afterSnapshot, err2 := e.takeWorkspaceSnapshot()
	if err2 != nil {
		log.Printf("Warning: failed to create after snapshot: %v", err2)
	} else {
		// Generate and save the report
		report := e.generateWorkspaceReport(beforeSnapshot, afterSnapshot)
		if reportErr := e.saveWorkspaceReport(report); reportErr != nil {
			log.Printf("Warning: failed to save workspace report: %v", reportErr)
		} else {
			log.Println("Workspace report saved to workspace-report.json")
		}
	}

	return err
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

// takeWorkspaceSnapshot creates a snapshot of the current workspace state
func (e *Engine) takeWorkspaceSnapshot() (WorkspaceSnapshot, error) {
	snapshot := WorkspaceSnapshot{
		Timestamp: time.Now(),
		Files:     make(map[string]FileInfo),
	}

	err := filepath.Walk(e.config.WorkspaceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get relative path from workspace root
		relPath, err := filepath.Rel(e.config.WorkspaceDir, path)
		if err != nil {
			return err
		}

		// Skip hidden files and directories
		if strings.HasPrefix(filepath.Base(relPath), ".") {
			return nil
		}

		fileInfo := FileInfo{
			Path:    relPath,
			Size:    info.Size(),
			ModTime: info.ModTime(),
			IsDir:   info.IsDir(),
		}

		// Calculate hash for files (not directories)
		if !info.IsDir() {
			hash, err := e.calculateFileHash(path)
			if err != nil {
				log.Printf("Warning: failed to hash file %s: %v", relPath, err)
				hash = ""
			}
			fileInfo.Hash = hash
		}

		snapshot.Files[relPath] = fileInfo
		return nil
	})

	return snapshot, err
}

// calculateFileHash computes MD5 hash of a file
func (e *Engine) calculateFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// generateWorkspaceReport compares two snapshots and generates a detailed report
func (e *Engine) generateWorkspaceReport(before, after WorkspaceSnapshot) WorkspaceReport {
	report := WorkspaceReport{
		Before: before,
		After:  after,
		Added:  []string{},
		Removed: []string{},
		Modified: []string{},
	}

	// Find added files
	for path := range after.Files {
		if _, exists := before.Files[path]; !exists {
			report.Added = append(report.Added, path)
		}
	}

	// Find removed files
	for path := range before.Files {
		if _, exists := after.Files[path]; !exists {
			report.Removed = append(report.Removed, path)
		}
	}

	// Find modified files
	for path, afterFile := range after.Files {
		if beforeFile, exists := before.Files[path]; exists {
			// Check if file was modified (different hash, size, or mod time)
			if !afterFile.IsDir && !beforeFile.IsDir {
				if afterFile.Hash != beforeFile.Hash {
					report.Modified = append(report.Modified, path)
				}
			}
		}
	}

	// Sort for consistent output
	sort.Strings(report.Added)
	sort.Strings(report.Removed)
	sort.Strings(report.Modified)

	// Generate summary
	report.Summary = e.generateSummary(report)

	return report
}

// generateSummary creates a human-readable summary of changes
func (e *Engine) generateSummary(report WorkspaceReport) string {
	var summary strings.Builder
	
	summary.WriteString(fmt.Sprintf("Workspace changes from %s to %s:\n",
		report.Before.Timestamp.Format("2006-01-02 15:04:05"),
		report.After.Timestamp.Format("2006-01-02 15:04:05")))
	
	summary.WriteString(fmt.Sprintf("- Files added: %d\n", len(report.Added)))
	summary.WriteString(fmt.Sprintf("- Files removed: %d\n", len(report.Removed)))
	summary.WriteString(fmt.Sprintf("- Files modified: %d\n", len(report.Modified)))
	
	if len(report.Added) > 0 {
		summary.WriteString("\nAdded files:\n")
		for _, file := range report.Added {
			summary.WriteString(fmt.Sprintf("  + %s\n", file))
		}
	}
	
	if len(report.Removed) > 0 {
		summary.WriteString("\nRemoved files:\n")
		for _, file := range report.Removed {
			summary.WriteString(fmt.Sprintf("  - %s\n", file))
		}
	}
	
	if len(report.Modified) > 0 {
		summary.WriteString("\nModified files:\n")
		for _, file := range report.Modified {
			beforeInfo := report.Before.Files[file]
			afterInfo := report.After.Files[file]
			summary.WriteString(fmt.Sprintf("  ~ %s (size: %d->%d bytes)\n", 
				file, beforeInfo.Size, afterInfo.Size))
		}
	}
	
	return summary.String()
}

// saveWorkspaceReport saves the workspace report to a JSON file
func (e *Engine) saveWorkspaceReport(report WorkspaceReport) error {
	reportPath := filepath.Join(e.config.WorkspaceDir, "workspace-report.json")
	
	// Pretty print JSON
	jsonData, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal report: %v", err)
	}
	
	if err := os.WriteFile(reportPath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write report file: %v", err)
	}
	
	// Also save a human-readable summary
	summaryPath := filepath.Join(e.config.WorkspaceDir, "workspace-summary.txt")
	if err := os.WriteFile(summaryPath, []byte(report.Summary), 0644); err != nil {
		log.Printf("Warning: failed to write summary file: %v", err)
	}
	
	// Print summary to console
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("WORKSPACE CHANGE REPORT")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println(report.Summary)
	fmt.Println(strings.Repeat("=", 60))
	
	return nil
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