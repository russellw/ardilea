package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// FileBasedTester provides file-based testing for BASIC interpreters
type FileBasedTester struct {
	interpreterPath string
	testsDir        string
	expectedDir     string
	errorsDir       string
}

// NewFileBasedTester creates a new file-based tester
func NewFileBasedTester(interpreterPath string) *FileBasedTester {
	return &FileBasedTester{
		interpreterPath: interpreterPath,
		testsDir:        "tests/basic",
		expectedDir:     "tests/expected",
		errorsDir:       "tests/errors",
	}
}

// RunBasicFile executes a BASIC file and returns the output
func (fbt *FileBasedTester) RunBasicFile(filename string) (string, error) {
	cmd := exec.Command(fbt.interpreterPath, filename)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("interpreter error: %v, stderr: %s", err, stderr.String())
	}

	return stdout.String(), nil
}

// ReadExpectedOutput reads the expected output file
func (fbt *FileBasedTester) ReadExpectedOutput(testName string) (string, error) {
	expectedFile := filepath.Join(fbt.expectedDir, testName+".txt")
	content, err := ioutil.ReadFile(expectedFile)
	if err != nil {
		return "", fmt.Errorf("failed to read expected output %s: %v", expectedFile, err)
	}
	return string(content), nil
}

// GetBasicFiles returns all .bas files in the tests directory
func (fbt *FileBasedTester) GetBasicFiles() ([]string, error) {
	files, err := filepath.Glob(filepath.Join(fbt.testsDir, "*.bas"))
	if err != nil {
		return nil, err
	}
	return files, nil
}

// GetErrorFiles returns all .bas files in the errors directory
func (fbt *FileBasedTester) GetErrorFiles() ([]string, error) {
	files, err := filepath.Glob(filepath.Join(fbt.errorsDir, "*.bas"))
	if err != nil {
		return nil, err
	}
	return files, nil
}

// GetTestName extracts test name from file path
func (fbt *FileBasedTester) GetTestName(filePath string) string {
	base := filepath.Base(filePath)
	return strings.TrimSuffix(base, ".bas")
}

// TestBasicInterpreterFilesBased runs file-based integration tests
func TestBasicInterpreterFilesBased(t *testing.T) {
	interpreterPath := os.Getenv("BASIC_INTERPRETER")
	if interpreterPath == "" {
		t.Skip("BASIC_INTERPRETER environment variable not set")
	}

	tester := NewFileBasedTester(interpreterPath)

	// Get all test files
	testFiles, err := tester.GetBasicFiles()
	if err != nil {
		t.Fatalf("Failed to get test files: %v", err)
	}

	if len(testFiles) == 0 {
		t.Fatal("No test files found in tests/basic/")
	}

	// Run each test
	for _, testFile := range testFiles {
		testName := tester.GetTestName(testFile)
		t.Run(testName, func(t *testing.T) {
			// Run the BASIC program
			actualOutput, err := tester.RunBasicFile(testFile)
			if err != nil {
				t.Fatalf("Failed to run %s: %v", testFile, err)
			}

			// Read expected output
			expectedOutput, err := tester.ReadExpectedOutput(testName)
			if err != nil {
				t.Fatalf("Failed to read expected output for %s: %v", testName, err)
			}

			// Compare outputs
			if actualOutput != expectedOutput {
				t.Errorf("Output mismatch for %s\nExpected:\n%s\nActual:\n%s", 
					testName, expectedOutput, actualOutput)
			}
		})
	}
}

// TestBasicInterpreterErrors tests error conditions
func TestBasicInterpreterErrors(t *testing.T) {
	interpreterPath := os.Getenv("BASIC_INTERPRETER")
	if interpreterPath == "" {
		t.Skip("BASIC_INTERPRETER environment variable not set")
	}

	tester := NewFileBasedTester(interpreterPath)

	// Get all error test files
	errorFiles, err := tester.GetErrorFiles()
	if err != nil {
		t.Fatalf("Failed to get error test files: %v", err)
	}

	if len(errorFiles) == 0 {
		t.Skip("No error test files found in tests/errors/")
	}

	// Run each error test
	for _, errorFile := range errorFiles {
		testName := tester.GetTestName(errorFile)
		t.Run(testName, func(t *testing.T) {
			// This should fail
			_, err := tester.RunBasicFile(errorFile)
			if err == nil {
				t.Errorf("Expected %s to fail, but it succeeded", testName)
			}
		})
	}
}

// TestBasicInterpreterManualExamples provides some manual verification
func TestBasicInterpreterManualExamples(t *testing.T) {
	interpreterPath := os.Getenv("BASIC_INTERPRETER")
	if interpreterPath == "" {
		t.Skip("BASIC_INTERPRETER environment variable not set")
	}

	tester := NewFileBasedTester(interpreterPath)

	// Test sample program
	if _, err := os.Stat("test_sample.bas"); err == nil {
		t.Run("sample_program", func(t *testing.T) {
			output, err := tester.RunBasicFile("test_sample.bas")
			if err != nil {
				t.Fatalf("Sample program failed: %v", err)
			}
			
			// Basic sanity checks
			if !strings.Contains(output, "BASIC Interpreter Test") {
				t.Error("Sample program output doesn't contain expected header")
			}
			if !strings.Contains(output, "Program completed successfully") {
				t.Error("Sample program didn't complete successfully")
			}
		})
	}
}

// BenchmarkBasicInterpreter benchmarks the interpreter
func BenchmarkBasicInterpreter(b *testing.B) {
	interpreterPath := os.Getenv("BASIC_INTERPRETER")
	if interpreterPath == "" {
		b.Skip("BASIC_INTERPRETER environment variable not set")
	}

	tester := NewFileBasedTester(interpreterPath)

	// Use factorial test for benchmarking
	testFile := "tests/basic/factorial.bas"
	if _, err := os.Stat(testFile); err != nil {
		b.Skip("Factorial test file not found")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := tester.RunBasicFile(testFile)
		if err != nil {
			b.Fatalf("Benchmark failed: %v", err)
		}
	}
}