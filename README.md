<h1><picture><source media="(prefers-color-scheme: dark)" srcset="brand/logo-dark.svg"><img src="brand/logo.svg" width="32" height="32" alt="" align="absmiddle"></picture> K-sense</h1>

> **Status: idea / planned.** Work has not started. This repository holds the plan and a CLI skeleton only.

K-sense is a compact coding language model specialized for Kotlin. It is meant to be fine-tuned offline. The model does not rewrite its own weights at runtime: all weight changes happen in a separate training cycle.

The intended training result is a base model plus a LoRA adapter in a standard format (for example safetensors), served through a standard OpenAI-compatible endpoint (for example vLLM, llama.cpp, or Ollama), so any agent or framework can call it.

Base-model directories exist under `models/`, but the weight files have not been uploaded yet.

## Quick start

The CLI is a Go stub with three commands. After a successful flag parse they print a "not implemented yet" status. `collect` and `train` also check that the named model directory exists.

```
cd cli && go build -o k-sense .
./k-sense collect --model qwen2.5-coder-7b
./k-sense train --model qwen2.5-coder-7b
./k-sense eval
```

Known model names are `qwen2.5-coder-7b` (default) and `deepseek-coder-v2-lite`. Every flag is listed in [Architecture: Go CLI](docs/README.md#8-go-cli).

## Documentation

- [Architecture](docs/README.md): goals, training stack, data plan, CLI flags and model storage.
- [models/](models/README.md): Git LFS layout for base-model weights.

## License

Not declared in this repository yet.
