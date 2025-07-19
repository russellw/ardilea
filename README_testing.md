# Language-Agnostic BASIC Interpreter Testing

This test suite provides language-agnostic integration testing for BASIC interpreters. Instead of being tied to a specific implementation, it executes any BASIC interpreter as an external process.

## Usage

Set the `BASIC_INTERPRETER` environment variable to the path of your BASIC interpreter executable:

```bash
# Test with Go implementation
BASIC_INTERPRETER=./basic go test -v basic_test.go

# Test with any other BASIC interpreter
BASIC_INTERPRETER=/path/to/your/basic go test -v basic_test.go
```

## Interpreter Requirements

Your BASIC interpreter must:

1. Accept a filename as a command-line argument
2. Execute the BASIC program in that file
3. Print output to stdout
4. Exit with non-zero status on errors
5. Support line-numbered BASIC syntax

Example usage of your interpreter:
```bash
./your_basic_interpreter program.bas
```

## Test Coverage

The test suite covers:

- **Basic Operations**: PRINT, LET, arithmetic
- **Control Flow**: GOTO, IF-THEN statements  
- **Loops**: FOR-NEXT loops (including nested)
- **Variables**: Numeric and string variables
- **Line Numbers**: Proper ordering and gaps
- **Comments**: REM statements
- **Error Handling**: Invalid syntax, undefined line numbers
- **Complex Programs**: Factorial calculation

## Building the Go Reference Implementation

```bash
go build -o basic interpreter.go
```

## Running Tests

```bash
# Build and test the Go implementation
go build -o basic interpreter.go
BASIC_INTERPRETER=./basic go test -v basic_test.go

# Run benchmarks
BASIC_INTERPRETER=./basic go test -bench=. basic_test.go
```

## Adding New Tests

Add test cases to the `tests` slice in `TestBasicInterpreterIntegration()`:

```go
{
    name:     "Test description",
    program:  `10 PRINT "Your BASIC program"`,
    expected: []string{"Expected output line 1", "Expected output line 2"},
},
```

## Error Tests

Error test cases should set `wantErr: true` and will pass if the interpreter exits with a non-zero status.

This approach allows testing any BASIC interpreter implementation (Python, C, Rust, etc.) with the same comprehensive test suite.