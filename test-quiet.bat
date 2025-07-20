@echo off
REM BASIC Interpreter Test Runner for Windows (Quiet Mode)
REM This batch file builds the reference implementation and runs tests without verbose output

echo Building BASIC interpreter...
go build -o basic.exe basic_reference_impl.go
if errorlevel 1 (
    echo ERROR: Failed to build BASIC interpreter
    exit /b 1
)

echo.
echo Running tests...
echo.
go run test_runner.go basic.exe
if errorlevel 1 (
    echo.
    echo Some tests failed. Check output above.
    exit /b 1
)

echo.
echo All tests passed successfully!