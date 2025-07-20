package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// BasicTester provides file-based testing for BASIC interpreters
type BasicTester struct {
	interpreterPath string
	testsDir        string
	expectedDir     string
	errorsDir       string
	passCount       int
	failCount       int
	verbose         bool
}

// NewBasicTester creates a new file-based tester
func NewBasicTester(interpreterPath string, verbose bool) *BasicTester {
	return &BasicTester{
		interpreterPath: interpreterPath,
		testsDir:        "tests/basic",
		expectedDir:     "tests/expected",
		errorsDir:       "tests/errors",
		passCount:       0,
		failCount:       0,
		verbose:         verbose,
	}
}

// RunBasicFile executes a BASIC file and returns the output
func (bt *BasicTester) RunBasicFile(filename string) (string, error) {
	cmd := exec.Command(bt.interpreterPath, filename)
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
func (bt *BasicTester) ReadExpectedOutput(testName string) (string, error) {
	expectedFile := filepath.Join(bt.expectedDir, testName+".txt")
	content, err := ioutil.ReadFile(expectedFile)
	if err != nil {
		return "", fmt.Errorf("failed to read expected output %s: %v", expectedFile, err)
	}
	return string(content), nil
}

// GetBasicFiles returns all .bas files in the tests directory
func (bt *BasicTester) GetBasicFiles() ([]string, error) {
	files, err := filepath.Glob(filepath.Join(bt.testsDir, "*.bas"))
	if err != nil {
		return nil, err
	}
	return files, nil
}

// GetErrorFiles returns all .bas files in the errors directory
func (bt *BasicTester) GetErrorFiles() ([]string, error) {
	files, err := filepath.Glob(filepath.Join(bt.errorsDir, "*.bas"))
	if err != nil {
		return nil, err
	}
	return files, nil
}

// GetTestName extracts test name from file path
func (bt *BasicTester) GetTestName(filePath string) string {
	base := filepath.Base(filePath)
	return strings.TrimSuffix(base, ".bas")
}

// RunSuccessTests runs all success tests and reports results
func (bt *BasicTester) RunSuccessTests() {
	fmt.Println("=== Running Success Tests ===")
	
	testFiles, err := bt.GetBasicFiles()
	if err != nil {
		fmt.Printf("Error getting test files: %v\n", err)
		return
	}

	if len(testFiles) == 0 {
		fmt.Println("No test files found in tests/basic/")
		return
	}

	for _, testFile := range testFiles {
		testName := bt.GetTestName(testFile)
		fmt.Printf("Running %s... ", testName)

		// Read BASIC source code for verbose output
		var sourceCode string
		if bt.verbose {
			if content, err := ioutil.ReadFile(testFile); err == nil {
				sourceCode = strings.TrimSpace(string(content))
			}
		}

		// Run the BASIC program
		actualOutput, err := bt.RunBasicFile(testFile)
		if err != nil {
			fmt.Printf("FAIL (execution error: %v)\n", err)
			if bt.verbose && sourceCode != "" {
				fmt.Printf("  BASIC code:\n%s\n", bt.indentLines(sourceCode))
			}
			bt.failCount++
			continue
		}

		// Read expected output
		expectedOutput, err := bt.ReadExpectedOutput(testName)
		if err != nil {
			fmt.Printf("FAIL (missing expected output: %v)\n", err)
			if bt.verbose && sourceCode != "" {
				fmt.Printf("  BASIC code:\n%s\n", bt.indentLines(sourceCode))
			}
			bt.failCount++
			continue
		}

		// Compare outputs
		if actualOutput == expectedOutput {
			fmt.Println("PASS")
			if bt.verbose {
				if sourceCode != "" {
					fmt.Printf("  BASIC code:\n%s\n", bt.indentLines(sourceCode))
				}
				fmt.Printf("  Output: %q\n", actualOutput)
			}
			bt.passCount++
		} else {
			fmt.Printf("FAIL (output mismatch)\n")
			if bt.verbose && sourceCode != "" {
				fmt.Printf("  BASIC code:\n%s\n", bt.indentLines(sourceCode))
			}
			fmt.Printf("  Expected: %q\n", expectedOutput)
			fmt.Printf("  Actual:   %q\n", actualOutput)
			bt.failCount++
		}
	}
}

// RunErrorTests runs all error tests and reports results
func (bt *BasicTester) RunErrorTests() {
	fmt.Println("\n=== Running Error Tests ===")
	
	errorFiles, err := bt.GetErrorFiles()
	if err != nil {
		fmt.Printf("Error getting error test files: %v\n", err)
		return
	}

	if len(errorFiles) == 0 {
		fmt.Println("No error test files found in tests/errors/")
		return
	}

	for _, errorFile := range errorFiles {
		testName := bt.GetTestName(errorFile)
		fmt.Printf("Running %s... ", testName)

		// Read BASIC source code for verbose output
		var sourceCode string
		if bt.verbose {
			if content, err := ioutil.ReadFile(errorFile); err == nil {
				sourceCode = strings.TrimSpace(string(content))
			}
		}

		// This should fail
		output, err := bt.RunBasicFile(errorFile)
		if err != nil {
			fmt.Println("PASS (correctly failed)")
			if bt.verbose {
				if sourceCode != "" {
					fmt.Printf("  BASIC code:\n%s\n", bt.indentLines(sourceCode))
				}
				fmt.Printf("  Error: %v\n", err)
			}
			bt.passCount++
		} else {
			fmt.Println("FAIL (should have failed but succeeded)")
			if bt.verbose {
				if sourceCode != "" {
					fmt.Printf("  BASIC code:\n%s\n", bt.indentLines(sourceCode))
				}
				fmt.Printf("  Unexpected output: %q\n", output)
			}
			bt.failCount++
		}
	}
}

// RunManualTests runs some manual verification tests
func (bt *BasicTester) RunManualTests() {
	fmt.Println("\n=== Running Manual Tests ===")
	
	// Test sample program if it exists
	if _, err := os.Stat("test_sample.bas"); err == nil {
		fmt.Printf("Running test_sample.bas... ")
		output, err := bt.RunBasicFile("test_sample.bas")
		if err != nil {
			fmt.Printf("FAIL (execution error: %v)\n", err)
			bt.failCount++
		} else {
			// Basic sanity checks
			if strings.Contains(output, "BASIC Interpreter Test") && 
			   strings.Contains(output, "Program completed successfully") {
				fmt.Println("PASS")
				if bt.verbose {
					fmt.Printf("  Output: %q\n", output)
				}
				bt.passCount++
			} else {
				fmt.Println("FAIL (unexpected output)")
				if bt.verbose {
					fmt.Printf("  Output: %q\n", output)
				}
				bt.failCount++
			}
		}
	}
}

// PrintSummary prints the test results summary
func (bt *BasicTester) PrintSummary() {
	fmt.Println("\n=== Test Summary ===")
	total := bt.passCount + bt.failCount
	fmt.Printf("Tests run: %d\n", total)
	fmt.Printf("Passed: %d\n", bt.passCount)
	fmt.Printf("Failed: %d\n", bt.failCount)
	
	if bt.failCount == 0 {
		fmt.Println("✅ All tests passed!")
	} else {
		fmt.Printf("❌ %d test(s) failed\n", bt.failCount)
	}
}

// indentLines adds 4-space indentation to each line
func (bt *BasicTester) indentLines(text string) string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = "    " + line
	}
	return strings.Join(lines, "\n")
}

// HasFailures returns true if any tests failed
func (bt *BasicTester) HasFailures() bool {
	return bt.failCount > 0
}

func main() {
	var interpreterPath string
	var verbose bool
	
	// Parse command line arguments
	args := os.Args[1:]
	for _, arg := range args {
		if arg == "-v" || arg == "--verbose" {
			verbose = true
		} else if !strings.HasPrefix(arg, "-") {
			interpreterPath = arg
			break
		}
	}
	
	// Fall back to environment variable if no interpreter specified
	if interpreterPath == "" {
		interpreterPath = os.Getenv("BASIC_INTERPRETER")
	}
	
	if interpreterPath == "" {
		fmt.Println("Usage:")
		fmt.Println("  go run test_runner.go [options] <interpreter_executable>")
		fmt.Println("  or")
		fmt.Println("  BASIC_INTERPRETER=./basic go run test_runner.go [options]")
		fmt.Println()
		fmt.Println("Options:")
		fmt.Println("  -v, --verbose    Show detailed output for each test")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  go run test_runner.go ./basic")
		fmt.Println("  go run test_runner.go -v ./basic")
		fmt.Println("  go run test_runner.go --verbose /usr/local/bin/my_basic")
		os.Exit(1)
	}

	// Check if interpreter exists
	if _, err := os.Stat(interpreterPath); os.IsNotExist(err) {
		fmt.Printf("Error: Interpreter not found at %s\n", interpreterPath)
		os.Exit(1)
	}
	
	// Fix relative path issue - if path doesn't start with ./ or /, prepend ./
	if !strings.HasPrefix(interpreterPath, "/") && !strings.HasPrefix(interpreterPath, "./") && !strings.HasPrefix(interpreterPath, "../") {
		interpreterPath = "./" + interpreterPath
	}

	fmt.Printf("Testing BASIC interpreter: %s\n", interpreterPath)
	if verbose {
		fmt.Println("Verbose mode enabled - showing detailed output")
	}
	
	tester := NewBasicTester(interpreterPath, verbose)
	
	// Run all test suites
	tester.RunSuccessTests()
	tester.RunErrorTests()
	tester.RunManualTests()
	
	// Print summary and exit with appropriate code
	tester.PrintSummary()
	
	if tester.HasFailures() {
		os.Exit(1)
	}
}