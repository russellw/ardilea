@echo off
REM Advanced LLM Response Time Test
REM This tests the Ollama server with complex programming prompts

echo Advanced LLM Response Time Test
echo ==================================

REM Check if Go is available
go version >nul 2>&1
if errorlevel 1 (
    echo ERROR: Go is not installed or not in PATH
    echo Please install Go to run this test
    exit /b 1
)

echo Go is available, running advanced LLM test...
echo.
echo WARNING: This test uses complex prompts that may take 2-10 minutes each.
echo The complete test could take 30-60 minutes to finish.
echo.

set /p CONTINUE="Continue with advanced test? (y/N): "
if /i not "%CONTINUE%"=="y" (
    echo Test cancelled.
    exit /b 0
)

echo.
echo Starting advanced test...
REM Run the test with timestamps
echo Test started at %TIME%
go run test_llm_advanced.go
echo Test completed at %TIME%

echo.
echo Analysis:
echo - Advanced prompts should take 2-10 minutes each for complex tasks
echo - Simple responses (under 30 seconds) may indicate incomplete processing
echo - Timeouts on advanced prompts may indicate server resource limits
echo - Compare with simple test results to identify complexity handling
echo.
pause