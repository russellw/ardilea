@echo off
REM Ardilea LLM Engine Runner for Windows
REM This batch file builds and runs the containerized LLM engine

echo Ardilea LLM Engine - Windows Runner
echo ===================================

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

REM Create workspace directory if it doesn't exist
if not exist workspace mkdir workspace

REM Build the container
echo Building Ardilea Engine container...
docker build -t ardilea-engine .
if errorlevel 1 (
    echo ERROR: Failed to build container
    exit /b 1
)

echo.
echo Container built successfully!
echo.

REM Check if config.json exists, if not create default
if not exist config.json (
    echo Creating default config.json...
    echo {> config.json
    echo   "ollama_server": "192.168.0.63:11434",>> config.json
    echo   "model_name": "qwen3:30b",>> config.json
    echo   "workspace_dir": "/workspace">> config.json
    echo }>> config.json
    echo Default configuration created.
    echo.
)

REM Show current configuration
echo Current configuration:
type config.json
echo.

REM Run the container
echo Starting Ardilea Engine...
echo.
echo Note: The engine will attempt to connect to Ollama server at the configured address.
echo Make sure your Ollama server is running and accessible.
echo.
echo Press Ctrl+C to stop the engine.
echo.

docker run -it --rm ^
    --name ardilea-engine ^
    -v "%cd%\workspace:/workspace/output" ^
    -v "%cd%\config.json:/workspace/config.json:ro" ^
    ardilea-engine

echo.
echo Engine stopped.