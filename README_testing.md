# File-Based BASIC Interpreter Testing

This test suite provides file-based, language-agnostic integration testing for BASIC interpreters. Test programs are stored as separate `.bas` files with corresponding expected output files, making the specification crystal clear and easy to understand.

## Directory Structure

```
tests/
├── basic/              # BASIC test programs
│   ├── hello.bas
│   ├── arithmetic.bas
│   ├── for_loop.bas
│   └── ...
├── expected/           # Expected output files
│   ├── hello.txt
│   ├── arithmetic.txt
│   ├── for_loop.txt
│   └── ...
└── errors/             # Programs that should fail
    ├── invalid_goto.bas
    ├── syntax_error.bas
    └── ...
```

## Usage

Set the `BASIC_INTERPRETER` environment variable to the path of your BASIC interpreter executable:

```bash
# Test with Go implementation
BASIC_INTERPRETER=./basic go run test_runner.go

# Test with any other BASIC interpreter
BASIC_INTERPRETER=/path/to/your/basic go run test_runner.go
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
BASIC_INTERPRETER=./basic go run test_runner.go

# Test individual programs manually
./basic tests/basic/hello.bas
./basic tests/basic/factorial.bas
```

## Adding New Tests

1. **Create a test program**: Add a new `.bas` file in `tests/basic/`
   ```basic
   10 PRINT "Hello, Test!"
   20 LET A = 42
   30 PRINT A
   ```

2. **Generate expected output**: Run the program and save output
   ```bash
   ./basic tests/basic/mytest.bas > tests/expected/mytest.txt
   ```

3. **The test runner automatically discovers and runs the new test**

## Error Tests

Programs in `tests/errors/` are expected to fail and will pass the test if the interpreter exits with a non-zero status.

To add error tests, simply add `.bas` files to `tests/errors/` - no expected output files needed.

## Benefits of File-Based Testing

- **Clear Specification**: Each `.bas` file clearly shows what features need implementing
- **Easy Debugging**: Run individual programs manually to debug issues
- **Language Agnostic**: Test any BASIC interpreter implementation (Python, C, Rust, etc.)
- **Version Control Friendly**: Git diffs show actual BASIC code changes
- **Self-Documenting**: Filename and program content describe the test
- **No Escaping**: BASIC programs are written in pure BASIC syntax with highlighting
- **No Dependencies**: Simple `go run` execution, no testing framework needed
- **Standalone**: Works without go.mod or module system

## Test Runner Output

The test runner provides clear output showing pass/fail for each test:

```
=== Running Success Tests ===
Running hello... PASS
Running arithmetic... PASS
Running for_loop... PASS
...

=== Running Error Tests ===
Running invalid_goto... PASS (correctly failed)
...

=== Test Summary ===
Tests run: 19
Passed: 19
Failed: 0
✅ All tests passed!
```

This approach makes the BASIC language specification crystal clear and easy to understand for anyone implementing a BASIC interpreter.