# models/

Веса моделей хранятся через **Git LFS**, а не обычными git-объектами. Шаблоны файлов (`*.safetensors`, `*.bin`, `*.gguf`, `*.pt` внутри `models/`) заданы в `/.gitattributes`.

| Каталог | Модель | Роль |
| --- | --- | --- |
| `qwen2.5-coder-7b/` | [Qwen/Qwen2.5-Coder-7B](https://huggingface.co/Qwen/Qwen2.5-Coder-7B) | основная |
| `deepseek-coder-v2-lite/` | [deepseek-ai/DeepSeek-Coder-V2-Lite-Base](https://huggingface.co/deepseek-ai/DeepSeek-Coder-V2-Lite-Base) | лёгкая альтернатива |

При `git clone` (если установлен `git lfs`) веса скачиваются автоматически. Если клон сделан без LFS, их можно докачать:

```
git lfs install
git lfs pull --include "models/qwen2.5-coder-7b/**"
```

**Статус:** сами веса в репозиторий пока не загружены. Сейчас в каталогах лежат только эти описания.
