package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// knownModels maps model names to their directories under models/.
// Weights in those directories are tracked with Git LFS.
var knownModels = map[string]string{
	"qwen2.5-coder-7b":       "qwen2.5-coder-7b",
	"deepseek-coder-v2-lite": "deepseek-coder-v2-lite",
}

const defaultModel = "qwen2.5-coder-7b"

var lfsPointer = []byte("version https://git-lfs.github.com/spec/v1")

// resolveModel returns the local path of a model under modelsDir and
// reports whether its weight files are present (not LFS pointers).
func resolveModel(modelsDir, name string) (string, bool, error) {
	dir, ok := knownModels[name]
	if !ok {
		return "", false, fmt.Errorf("unknown model %q", name)
	}
	path := filepath.Join(modelsDir, dir)
	if _, err := os.Stat(path); err != nil {
		return "", false, fmt.Errorf("model directory %s: %w", path, err)
	}
	entries, err := os.ReadDir(path)
	if err != nil {
		return "", false, err
	}
	present := false
	for _, e := range entries {
		if e.IsDir() || !isWeightFile(e.Name()) {
			continue
		}
		head := make([]byte, len(lfsPointer))
		f, err := os.Open(filepath.Join(path, e.Name()))
		if err != nil {
			return "", false, err
		}
		n, _ := f.Read(head)
		f.Close()
		if bytes.Equal(head[:n], lfsPointer) {
			return path, false, nil
		}
		present = true
	}
	return path, present, nil
}

func isWeightFile(name string) bool {
	for _, ext := range []string{".safetensors", ".bin", ".gguf", ".pt"} {
		if strings.HasSuffix(name, ext) {
			return true
		}
	}
	return false
}

func modelStatus(modelsDir, name string) (string, error) {
	path, present, err := resolveModel(modelsDir, name)
	if err != nil {
		return "", err
	}
	if !present {
		fmt.Fprintf(os.Stderr, "warning: no weights in %s; run `git lfs pull --include \"%s/**\"`\n", path, path)
	}
	return path, nil
}
