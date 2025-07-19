package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// BasicInterpreterTester provides language-agnostic testing for BASIC interpreters
type BasicInterpreterTester struct {
	interpreterPath string
}

// NewBasicInterpreterTester creates a new tester for the given interpreter executable
func NewBasicInterpreterTester(interpreterPath string) *BasicInterpreterTester {
	return &BasicInterpreterTester{
		interpreterPath: interpreterPath,
	}
}

// RunProgram executes a BASIC program using the configured interpreter
func (bit *BasicInterpreterTester) RunProgram(program string) ([]string, error) {
	// Create temporary file for the program
	tmpFile, err := ioutil.TempFile("", "basic_program_*.bas")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write program to file
	if _, err := tmpFile.WriteString(program); err != nil {
		return nil, fmt.Errorf("failed to write program: %v", err)
	}
	tmpFile.Close()

	// Execute interpreter
	cmd := exec.Command(bit.interpreterPath, tmpFile.Name())
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("interpreter error: %v, stderr: %s", err, stderr.String())
	}

	// Parse output lines
	output := stdout.String()
	if output == "" {
		return []string{}, nil
	}

	lines := strings.Split(strings.TrimSuffix(output, "\n"), "\n")
	return lines, nil
}

// TestBasicInterpreterIntegration runs comprehensive integration tests
func TestBasicInterpreterIntegration(t *testing.T) {
	interpreterPath := os.Getenv("BASIC_INTERPRETER")
	if interpreterPath == "" {
		t.Skip("BASIC_INTERPRETER environment variable not set")
	}

	tester := NewBasicInterpreterTester(interpreterPath)

	tests := []struct {
		name     string
		program  string
		expected []string
	}{
		{
			name:     "Simple print statement",
			program:  `10 PRINT "Hello, World!"`,
			expected: []string{"Hello, World!"},
		},
		{
			name: "Multiple line program",
			program: `10 PRINT "First line"
20 PRINT "Second line"
30 PRINT "Third line"`,
			expected: []string{"First line", "Second line", "Third line"},
		},
		{
			name: "Line number ordering",
			program: `30 PRINT "Third"
10 PRINT "First"
20 PRINT "Second"`,
			expected: []string{"First", "Second", "Third"},
		},
		{
			name: "Variable assignment and usage",
			program: `10 LET A = 42
20 PRINT A`,
			expected: []string{"42"},
		},
		{
			name: "Arithmetic operations",
			program: `10 LET A = 10
20 LET B = 5
30 PRINT A + B
40 PRINT A - B
50 PRINT A * B
60 PRINT A / B`,
			expected: []string{"15", "5", "50", "2"},
		},
		{
			name: "GOTO statement",
			program: `10 PRINT "First"
20 GOTO 40
30 PRINT "This should not print"
40 PRINT "Last"`,
			expected: []string{"First", "Last"},
		},
		{
			name: "IF statement",
			program: `10 LET A = 10
20 IF A > 5 THEN PRINT "A is greater than 5"
30 IF A < 5 THEN PRINT "A is less than 5"
40 PRINT "Done"`,
			expected: []string{"A is greater than 5", "Done"},
		},
		{
			name: "FOR loop",
			program: `10 FOR I = 1 TO 3
20 PRINT I
30 NEXT I`,
			expected: []string{"1", "2", "3"},
		},
		{
			name: "Nested FOR loops",
			program: `10 FOR I = 1 TO 2
20 FOR J = 1 TO 2
30 PRINT I; J
40 NEXT J
50 NEXT I`,
			expected: []string{"1 1", "1 2", "2 1", "2 2"},
		},
		{
			name: "String operations",
			program: `10 LET A$ = "Hello"
20 LET B$ = "World"
30 PRINT A$; " "; B$; "!"`,
			expected: []string{"Hello   World !"},
		},
		{
			name: "Line number gaps",
			program: `100 PRINT "Line 100"
500 PRINT "Line 500"
1000 PRINT "Line 1000"`,
			expected: []string{"Line 100", "Line 500", "Line 1000"},
		},
		{
			name: "Program with comments",
			program: `10 REM This is a comment
20 PRINT "This should print"
30 REM Another comment
40 PRINT "This should also print"`,
			expected: []string{"This should print", "This should also print"},
		},
		{
			name: "END statement",
			program: `10 PRINT "Before END"
20 END
30 PRINT "After END - should not print"`,
			expected: []string{"Before END"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := tester.RunProgram(tt.program)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if len(output) != len(tt.expected) {
				t.Fatalf("Expected %d lines of output, got %d. Output: %v", len(tt.expected), len(output), output)
			}

			for i, expected := range tt.expected {
				if output[i] != expected {
					t.Errorf("Line %d: expected %q, got %q", i+1, expected, output[i])
				}
			}
		})
	}
}

// TestBasicInterpreterErrors tests error handling
func TestBasicInterpreterErrors(t *testing.T) {
	interpreterPath := os.Getenv("BASIC_INTERPRETER")
	if interpreterPath == "" {
		t.Skip("BASIC_INTERPRETER environment variable not set")
	}

	tester := NewBasicInterpreterTester(interpreterPath)

	errorTests := []struct {
		name    string
		program string
		wantErr bool
	}{
		{
			name:    "Invalid line number in GOTO",
			program: `10 PRINT "Start"
20 GOTO 999
30 PRINT "End"`,
			wantErr: true,
		},
		{
			name:    "Syntax error",
			program: `10 PRINT "Valid line"
20 INVALID_COMMAND
30 PRINT "Another valid line"`,
			wantErr: true,
		},
		{
			name:    "Division by zero",
			program: `10 LET A = 10
20 LET B = 0
30 PRINT A / B`,
			wantErr: true,
		},
		{
			name:    "NEXT without FOR",
			program: `10 PRINT "Start"
20 NEXT I
30 PRINT "End"`,
			wantErr: true,
		},
	}

	for _, tt := range errorTests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tester.RunProgram(tt.program)
			if (err != nil) != tt.wantErr {
				t.Errorf("Expected error: %v, got error: %v", tt.wantErr, err)
			}
		})
	}
}

// TestComplexProgram tests a complex factorial calculation
func TestComplexProgram(t *testing.T) {
	interpreterPath := os.Getenv("BASIC_INTERPRETER")
	if interpreterPath == "" {
		t.Skip("BASIC_INTERPRETER environment variable not set")
	}

	tester := NewBasicInterpreterTester(interpreterPath)

	program := `10 LET N = 5
20 LET F = 1
30 FOR I = 1 TO N
40 LET F = F * I
50 NEXT I
60 PRINT "Factorial of"; N; "is"; F
70 END`

	output, err := tester.RunProgram(program)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(output) != 1 {
		t.Fatalf("Expected 1 line of output, got %d", len(output))
	}

	expected := "Factorial of 5 is 120"
	if output[0] != expected {
		t.Errorf("Expected %q, got %q", expected, output[0])
	}
}

// Benchmark tests
func BenchmarkBasicInterpreter(b *testing.B) {
	interpreterPath := os.Getenv("BASIC_INTERPRETER")
	if interpreterPath == "" {
		b.Skip("BASIC_INTERPRETER environment variable not set")
	}

	tester := NewBasicInterpreterTester(interpreterPath)
	program := `10 FOR I = 1 TO 100
20 LET A = I * 2
30 NEXT I
40 PRINT "Done"`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := tester.RunProgram(program)
		if err != nil {
			b.Fatalf("Benchmark failed: %v", err)
		}
	}
}