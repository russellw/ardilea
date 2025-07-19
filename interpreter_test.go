package main

import (
	"strings"
	"testing"
)

func TestBasicInterpreterIntegration(t *testing.T) {
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
			interpreter := NewBasicInterpreter()
			err := interpreter.Run(tt.program)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			output := interpreter.GetOutput()
			if len(output) != len(tt.expected) {
				t.Fatalf("Expected %d lines of output, got %d", len(tt.expected), len(output))
			}

			for i, expected := range tt.expected {
				if output[i] != expected {
					t.Errorf("Line %d: expected %q, got %q", i+1, expected, output[i])
				}
			}
		})
	}
}

func TestBasicInterpreterErrors(t *testing.T) {
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
			interpreter := NewBasicInterpreter()
			err := interpreter.Run(tt.program)
			if (err != nil) != tt.wantErr {
				t.Errorf("Expected error: %v, got error: %v", tt.wantErr, err)
			}
		})
	}
}

func TestComplexProgram(t *testing.T) {
	// Test a factorial calculation program
	program := `10 LET N = 5
20 LET F = 1
30 FOR I = 1 TO N
40 LET F = F * I
50 NEXT I
60 PRINT "Factorial of"; N; "is"; F
70 END`

	interpreter := NewBasicInterpreter()
	err := interpreter.Run(program)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	output := interpreter.GetOutput()
	if len(output) != 1 {
		t.Fatalf("Expected 1 line of output, got %d", len(output))
	}

	expected := "Factorial of 5 is 120"
	if output[0] != expected {
		t.Errorf("Expected %q, got %q", expected, output[0])
	}
}

func TestProgramStateIsolation(t *testing.T) {
	interpreter := NewBasicInterpreter()

	// Run first program
	program1 := "10 LET A = 42"
	err := interpreter.Run(program1)
	if err != nil {
		t.Fatalf("Unexpected error in first program: %v", err)
	}

	// Run second program - should not have access to A from first program
	program2 := "10 PRINT A"
	err = interpreter.Run(program2)
	if err == nil {
		t.Fatal("Expected error when accessing undefined variable, but got none")
	}
	if !strings.Contains(err.Error(), "cannot evaluate expression") {
		t.Errorf("Expected 'cannot evaluate expression' error, got: %v", err)
	}
}

func TestLoadProgram(t *testing.T) {
	interpreter := NewBasicInterpreter()
	program := `10 PRINT "Line 1"
20 PRINT "Line 2"
30 PRINT "Line 3"`

	err := interpreter.LoadProgram(program)
	if err != nil {
		t.Fatalf("Unexpected error loading program: %v", err)
	}

	err = interpreter.Execute()
	if err != nil {
		t.Fatalf("Unexpected error executing program: %v", err)
	}

	output := interpreter.GetOutput()
	expected := []string{"Line 1", "Line 2", "Line 3"}

	if len(output) != len(expected) {
		t.Fatalf("Expected %d lines of output, got %d", len(expected), len(output))
	}

	for i, exp := range expected {
		if output[i] != exp {
			t.Errorf("Line %d: expected %q, got %q", i+1, exp, output[i])
		}
	}
}