package main

import (
	"context"
	"encoding/json"
	"firefly/alex/internal/app"
	"firefly/alex/internal/words"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	defaultURLsPath       = "./data/endg-urls"
	defaultWordsPath      = "./data/words.txt"
	defaultMinWordSize    = 3
	defaultTopWordsCount  = 10
	defaultTaskBufferSize = 300
	defaulWorkers         = 20
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	urlsPath := flag.String("urls-path", defaultURLsPath, "path of the file containing essay urls")
	wordsPath := flag.String("words-path", defaultWordsPath, "path of the file containing a list of words")
	minWordSize := flag.Int("min-word-size", defaultMinWordSize, "minimum word size")
	topWordsCount := flag.Int("top-words-count", defaultTopWordsCount, "number of most recurrent words to return")
	bufferSize := flag.Int("tasks-buffer-size", defaultTaskBufferSize, "max number of tasks loaded in the pipeline")
	workerCount := flag.Int("worker-count", defaulWorkers, "number of workers")
	flag.Parse()

	logger.Info("Loading word bank...")

	start := time.Now()
	wordBank, err := words.NewBank(*wordsPath, words.IsMinimumSize(*minWordSize), words.IsAlpha())
	if err != nil {
		logger.Error("Failed to load word bank", "path", *wordsPath, "error", err)
		return
	}

	logger.Info("Successfully loaded word bank", "duration", time.Since(start))

	logger.Info("Processing essays...")

	start = time.Now()
	essayWordCount := app.NewEssayWordCount(logger, *urlsPath, wordBank)
	result, err := essayWordCount.Run(ctx, *topWordsCount, *bufferSize, *workerCount)
	if err != nil {
		logger.Error("Failed to run essay word count", "error", err)
		return
	}

	logger.Info("Successfully processed essays", "duration", time.Since(start))

	jsonBytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		logger.Error("Failed to marshal result to JSON", "error", err)
		return
	}

	fmt.Println(string(jsonBytes))
}
