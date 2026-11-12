# K-sense: Architecture Overview

> **Status: idea / planned.** Nothing beyond this plan and a CLI skeleton exists yet.

This document explains each component of K-sense: what it does and why it is there. The Russian [README](../README.md) is the primary project page; this is the English reference.

## 1. Project goal

**What:** A compact LLM for code, specialized in Kotlin. It is produced by offline fine-tuning of an existing open coding model.

**Why:** Top Kotlin results today come from large proprietary models used through agents. Claude Opus 4.7 and GPT 5.5 lead the [JetBrains Kotlin Benchmark](https://kotlinlang.org/benchmark/). K-sense does not try to beat them. It aims for a small model that runs locally and knows Kotlin codebases well. The model never rewrites itself at runtime: all weight changes happen in a separate, controlled training cycle, which keeps behaviour reproducible and reviewable.

## 2. Base models

**What:**
- **Primary:** Qwen 2.5 Coder 7B.
- **Lighter alternative:** DeepSeek Coder V2 Lite, a mixture-of-experts model with 16B total and 2.4B active parameters.

**Why:** Starting from a strong open coding model is far cheaper than training from scratch. In the Qwen 2.5 Coder technical report (Alibaba, 2024), Qwen 2.5 Coder 7B beats DeepSeek Coder V2 Lite on most coding benchmarks: HumanEval, MBPP, BigCodeBench, LiveCodeBench, Aider and CRUXEval. That makes it the default. DeepSeek Coder V2 Lite stays as an option for lower inference cost.

## 3. Training stack: Axolotl, Unsloth, LoRA

**What:**
- **Axolotl:** config-driven fine-tuning pipelines.
- **Unsloth:** fast, memory-efficient training kernels.
- **LoRA adapters:** small trainable weight deltas on top of the frozen base model.

**Why:** Axolotl gives reproducible, declarative training configs, and Unsloth makes 7B-class fine-tuning feasible on modest hardware. LoRA keeps the base weights untouched, so every training run produces a separate adapter. Rollback means switching back to a previous adapter or disabling it. Several adapters can be compared side by side.

## 4. Replay and EWC

**What:**
- **Replay:** every training batch mixes fresh data with chunks of older data.
- **Elastic Weight Consolidation (EWC):** an optional regularizer. It penalizes large changes to weights that were important for earlier tasks.

**Why:** Replay is the main mechanism for overwriting stale knowledge. Fresh examples replace outdated patterns, such as deprecated APIs or old library versions. The replayed older data keeps still-valid knowledge from being lost, so the model avoids catastrophic forgetting. EWC is an extra safety net, used when replay alone is not enough.

## 5. Training data

**What:** The model learns from the history of changes, not from whole repositories:
- diffs;
- commit messages;
- pull request descriptions;
- tests that failed and were later fixed;
- review comments.

Target mix: roughly **80% Kotlin, 20% other languages** such as JavaScript and Rust.

**Why:** Change history shows how code evolves and why: what was wrong, how it was fixed, and what reviewers expect. Whole-repo dumps mostly repeat static code the base model has already seen. The 20% share of other languages keeps general coding skills and cross-language context, for example Kotlin/JS or native interop, instead of narrowing the model too far.

## 6. Agent-agnostic design

**What:** The output is a standalone artifact: base weights plus a LoRA adapter in standard formats such as safetensors. It is served through a standard model endpoint, for example an OpenAI-compatible API via vLLM, llama.cpp or Ollama.

**Why:** The model is not tied to a single agent or framework. Any agent that can call a standard endpoint, or load standard weights, can use it. Swapping agents or upgrading the model never forces changes on the other side.

## 7. Separate staleness detection and training

**What:** Data collection and fine-tuning are separate commands, never combined in one step:
1. **collect** detects stale vs fresh data and builds a dataset with replay chunks.
2. **train** runs offline LoRA fine-tuning on an already collected dataset.
3. **eval** evaluates the adapter and optionally deploys it.

**Why:** Decoupling lets each step be run, repeated, inspected and rolled back on its own. A dataset can be reviewed before any GPU time is spent. Training can be re-run on the same dataset with different settings. A bad collection run never silently turns into a bad model.

## 8. Go CLI

**What:** A minimal CLI in `cli/`, currently with stub commands:

```
k-sense collect --model qwen2.5-coder-7b   # detect stale/fresh data, build dataset with replay
k-sense train   --model qwen2.5-coder-7b   # offline LoRA fine-tuning (optional --ewc)
k-sense eval    [--deploy]                 # evaluate the adapter; --deploy publishes it
```

Build it with `cd cli && go build -o k-sense .`

**Why Go:** It builds into a single static binary with no runtime dependency, and it compiles fast. That makes it easy to ship to any machine that does collection or training. The scope stays deliberately small: the CLI orchestrates the steps, while the heavy lifting stays in the training stack.

## 9. Model storage: Git LFS under `models/`

**What:**
- Base model weights live in `models/qwen2.5-coder-7b/` and `models/deepseek-coder-v2-lite/`.
- They are tracked with Git LFS. `.gitattributes` routes `*.safetensors`, `*.bin`, `*.gguf` and `*.pt` under `models/` through LFS.
- Cloning with `git lfs` installed pulls the weights. A clone made without LFS can fetch them with `git lfs pull`.
- `collect` and `train` resolve the model from the local `models/<name>` path and warn when only LFS pointers are present.

**Why:** Weights are versioned together with the code that uses them, so a commit always identifies the exact model files. Large binaries stay out of regular git objects, which keeps the repository history small. The weights themselves have not been uploaded yet. Note that GitHub LFS storage quotas are limited, and a 7B model is about 15 GB.
