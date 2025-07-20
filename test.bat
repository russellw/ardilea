@echo off
REM BASIC Interpreter Test Runner for Windows
REM This batch file builds the reference implementation and runs tests in verbose mode

echo Building BASIC interpreter...
go build -o basic.exe basic_reference_impl.go
if errorlevel 1 (
    echo ERROR: Failed to build BASIC interpreter
    pause
    exit /b 1
)

echo.
echo Running tests in verbose mode...
echo.
go run test_runner.go -v basic.exe
if errorlevel 1 (
    echo.
    echo Some tests failed. Check output above.
    pause
    exit /b 1
)

echo.
echo All tests passed successfully!
pause