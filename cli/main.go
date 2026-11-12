// Command k-sense is a minimal CLI for the K-sense project.
// Data collection, training, and evaluation are separate commands by design.
package main

import (
	"flag"
	"fmt"
	"os"
)

const usage = `k-sense: Kotlin-specialized coding LLM tooling (planned)

Usage:
  k-sense collect [flags]   detect stale vs fresh data and build a dataset with replay
  k-sense train   [flags]   run offline LoRA fine-tuning on a collected dataset
  k-sense eval    [flags]   evaluate an adapter; --deploy publishes it
`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}
	cmd, args := os.Args[1], os.Args[2:]
	var err error
	switch cmd {
	case "collect":
		err = runCollect(args)
	case "train":
		err = runTrain(args)
	case "eval":
		err = runEval(args)
	case "-h", "--help", "help":
		fmt.Print(usage)
		return
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n\n%s", cmd, usage)
		os.Exit(2)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func runCollect(args []string) error {
	fs := flag.NewFlagSet("collect", flag.ExitOnError)
	out := fs.String("out", "data/dataset", "output dataset directory")
	ratio := fs.Float64("kotlin-ratio", 0.8, "share of Kotlin samples")
	replay := fs.Float64("replay", 0.2, "share of old (replay) data per batch")
	model := fs.String("model", defaultModel, "model name under models/")
	modelsDir := fs.String("models-dir", "models", "LFS-tracked models directory")
	_ = fs.Parse(args)
	path, err := modelStatus(*modelsDir, *model)
	if err != nil {
		return err
	}
	fmt.Printf("collect: not implemented yet (model=%s out=%s kotlin-ratio=%.2f replay=%.2f)\n", path, *out, *ratio, *replay)
	return nil
}

func runTrain(args []string) error {
	fs := flag.NewFlagSet("train", flag.ExitOnError)
	data := fs.String("data", "data/dataset", "collected dataset directory")
	model := fs.String("model", defaultModel, "base model name under models/")
	modelsDir := fs.String("models-dir", "models", "LFS-tracked models directory")
	ewc := fs.Bool("ewc", false, "enable Elastic Weight Consolidation")
	_ = fs.Parse(args)
	path, err := modelStatus(*modelsDir, *model)
	if err != nil {
		return err
	}
	fmt.Printf("train: not implemented yet (data=%s model=%s ewc=%t)\n", *data, path, *ewc)
	return nil
}

func runEval(args []string) error {
	fs := flag.NewFlagSet("eval", flag.ExitOnError)
	adapter := fs.String("adapter", "out/adapter", "LoRA adapter directory")
	deploy := fs.Bool("deploy", false, "publish the adapter after evaluation")
	_ = fs.Parse(args)
	fmt.Printf("eval: not implemented yet (adapter=%s deploy=%t)\n", *adapter, *deploy)
	return nil
}
