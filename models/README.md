# models/

Model weights are stored with **Git LFS**, not as regular git objects. The file patterns (`*.safetensors`, `*.bin`, `*.gguf`, `*.pt` under `models/`) are set in `/.gitattributes`.

| Directory | Model | Role |
| --- | --- | --- |
| `qwen2.5-coder-7b/` | [Qwen/Qwen2.5-Coder-7B](https://huggingface.co/Qwen/Qwen2.5-Coder-7B) | primary |
| `deepseek-coder-v2-lite/` | [deepseek-ai/DeepSeek-Coder-V2-Lite-Base](https://huggingface.co/deepseek-ai/DeepSeek-Coder-V2-Lite-Base) | lighter alternative |

A `git clone` with `git lfs` installed fetches the weights automatically. If the clone was made without LFS, fetch them with:

```
git lfs install
git lfs pull --include "models/qwen2.5-coder-7b/**"
```

**Status:** the weights have not been uploaded yet. These directories currently hold only the descriptions.
