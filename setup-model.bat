@echo off
REM Setup LLM Model for Ardilea Engine
REM This batch file helps download and setup the required LLM model

echo Ardilea Model Setup - Windows
echo ===============================

REM Check if Docker is available
docker --version >nul 2>&1
if errorlevel 1 (
    echo ERROR: Docker is not installed or not in PATH
    echo Please install Docker Desktop and ensure it's running
    exit /b 1
)

REM Check if Docker is running
docker info >nul 2>&1
if errorlevel 1 (
    echo ERROR: Docker is not running
    echo Please start Docker Desktop and try again
    exit /b 1
)

echo Docker is available and running...
echo.

REM Read current config to get model name
set MODEL_NAME=qwen3:30b
if exist config.json (
    echo Reading model name from config.json...
    REM Note: This is a simplified approach - in production you'd parse JSON properly
    for /f "tokens=2 delims=:" %%a in ('findstr "model_name" config.json') do (
        set MODEL_LINE=%%a
    )
    REM Remove quotes and comma
    for /f "tokens=1 delims=," %%b in ("%MODEL_LINE%") do (
        set MODEL_NAME=%%b
    )
    set MODEL_NAME=!MODEL_NAME: =!
    set MODEL_NAME=!MODEL_NAME:"=!
)

echo Target model: %MODEL_NAME%
echo.

echo This will:
echo 1. Start Ollama server container
echo 2. Download the model: %MODEL_NAME%
echo 3. Keep the server running for the engine to use
echo.
echo Note: Model download may take significant time depending on model size and internet speed.
echo The qwen3:30b model is approximately 19GB.
echo.

set /p CONTINUE="Continue? (y/N): "
if /i not "%CONTINUE%"=="y" (
    echo Setup cancelled.
    exit /b 0
)

echo.
echo Starting Ollama server...
docker run -d --name ollama-server -p 11434:11434 -v ollama-data:/root/.ollama ollama/ollama:latest

REM Wait a moment for server to start
echo Waiting for Ollama server to start...
timeout /t 5 /nobreak >nul

echo.
echo Downloading model: %MODEL_NAME%
echo This may take a while...
docker exec ollama-server ollama pull %MODEL_NAME%

if errorlevel 1 (
    echo ERROR: Failed to download model
    echo Stopping Ollama server...
    docker stop ollama-server
    docker rm ollama-server
    exit /b 1
)

echo.
echo Model setup complete!
echo.
echo Ollama server is running at http://localhost:11434
echo You can now run the Ardilea engine with: run-engine.bat
echo.
echo To stop the Ollama server later, run:
echo   docker stop ollama-server
echo   docker rm ollama-server
echo.

echo Testing connection to Ollama server...
curl -s http://localhost:11434/api/tags >nul
if errorlevel 1 (
    echo WARNING: Could not connect to Ollama server
    echo Make sure the server is running and accessible
) else (
    echo SUCCESS: Ollama server is accessible
)

echo.
echo Setup complete!