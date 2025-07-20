@echo off
REM Ardilea Complete Stack Runner for Windows
REM This batch file runs the complete stack including Ollama server using Docker Compose

echo Ardilea Complete Stack - Windows Runner
echo ========================================

REM Check if Docker is available
docker --version >nul 2>&1
if errorlevel 1 (
    echo ERROR: Docker is not installed or not in PATH
    echo Please install Docker Desktop and ensure it's running
    exit /b 1
)

REM Check if Docker Compose is available
docker-compose --version >nul 2>&1
if errorlevel 1 (
    echo ERROR: Docker Compose is not available
    echo Please install Docker Desktop which includes Docker Compose
    exit /b 1
)

REM Check if Docker is running
docker info >nul 2>&1
if errorlevel 1 (
    echo ERROR: Docker is not running
    echo Please start Docker Desktop and try again
    exit /b 1
)

echo Docker and Docker Compose are available...
echo.

REM Create workspace directory if it doesn't exist
if not exist workspace mkdir workspace

REM Check if config.json exists, if not create default
if not exist config.json (
    echo Creating default config.json...
    echo {> config.json
    echo   "ollama_server": "ollama:11434",>> config.json
    echo   "model_name": "qwen3:30b",>> config.json
    echo   "workspace_dir": "/workspace">> config.json
    echo }>> config.json
    echo Default configuration created for containerized Ollama.
    echo.
)

REM Show current configuration
echo Current configuration:
type config.json
echo.

echo Starting complete Ardilea stack...
echo This includes:
echo - Ollama server (will download model if needed)
echo - Ardilea LLM Engine
echo.
echo Note: First run may take time to download the LLM model.
echo The Ollama server will be available at http://localhost:11434
echo.
echo Press Ctrl+C to stop the entire stack.
echo.

REM Start the stack
docker-compose up --build

echo.
echo Stack stopped.