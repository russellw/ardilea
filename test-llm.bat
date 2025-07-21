@echo off
REM LLM Response Time Test
REM This tests the Ollama server directly to measure response times

echo LLM Response Time Test
echo ======================

REM Check if Go is available
go version >nul 2>&1
if errorlevel 1 (
    echo ERROR: Go is not installed or not in PATH
    echo Please install Go to run this test
    exit /b 1
)

echo Go is available, running LLM test...
echo.

REM Run the test with timestamps
echo Test started at %TIME%
go run test_llm.go
echo Test completed at %TIME%

echo.
echo Analysis:
echo - Simple prompts should respond in 1-10 seconds
echo - Programming prompts may take 30-120 seconds  
echo - If any test hangs for over 5 minutes, there's likely a server issue
echo.
pause