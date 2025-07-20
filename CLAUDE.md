The goal of this project is to run an LLM in agent mode, trying to get it to write a Basic interpreter
The LLM will be running on an Ollama server on a local network
The IP address of the Ollama server should be read from a config file, default to 192.168.0.63
The name of the model should be read from a config file, default to qwen3:30b
The engine that calls the LLM should be written in Go
The engine and workspace should be in a Docker container

