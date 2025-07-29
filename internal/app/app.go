package app

import (
	"bufio"
	"context"
	"firefly/alex/internal/httpclient"
	"firefly/alex/internal/sources"
	"firefly/alex/internal/words"
	"fmt"
	"log/slog"
	"os"
	"sort"
	"strings"
	"sync"
	"unicode"
)

type (
	WordCount struct {
		Word  string
		Count uint64
	}
	EssayWordCount struct {
		urlsPath string
		wordBank *words.Bank

		logger     *slog.Logger
		httpClient *httpclient.RateLimitedRetryClient

		// concurrency:
		// - a single go routine produces tasks (loading essays) to taskChan
		// - each worker consumes tasks from taskChan and pushes its result to resultChan
		// - the main go routine consumes from resultChan and aggregates the results
		taskChan   chan string // TODO: could be extended to an interface to support different types of tasks, here contains URLs
		resultChan chan map[string]uint64
	}
)

func NewEssayWordCount(logger *slog.Logger, urlsPath string, wordBank *words.Bank) EssayWordCount {
	return EssayWordCount{
		urlsPath: urlsPath,
		wordBank: wordBank,
		logger:   logger.With(slog.String("context", "countWords")),
		httpClient: httpclient.NewRateLimitedRetryClient(
			100, // requests per second
			1,   // no burst
			"MyApp/0.1-alpha",
		),
	}
}

func (e *EssayWordCount) Run(ctx context.Context, maxTopWords, bufferSize, workerCount int) ([]WordCount, error) {
	e.taskChan = make(chan string, bufferSize)
	e.resultChan = make(chan map[string]uint64)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(e.taskChan)
		e.logger.Debug("urls loader started")

		if err := streamFileToChannel(ctx, e.urlsPath, e.taskChan); err != nil {
			e.logger.Error("failed to load url file", "path", e.urlsPath, "error", err)
		}

		e.logger.Debug("urls loader done")
	}()

	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go func(id int) {
			defer wg.Done()
			e.logger.Debug("worker %d started", "worker", id)

			e.processEssays(ctx, e.taskChan, e.resultChan)

			e.logger.Debug("worker %d done", "worker", id)
		}(i)
	}

	go func() {
		wg.Wait()
		close(e.resultChan)
	}()

	counts := map[string]uint64{}

	for result := range e.resultChan {
		counts = mergeCounts(counts, result)
		e.logger.Debug("result processed")
	}

	if ctx.Err() != nil {
		return nil, fmt.Errorf("work interrupted: %w", ctx.Err())
	}

	e.logger.Debug("all results processed")

	return topNWords(counts, maxTopWords), nil
}

func (e *EssayWordCount) processEssays(ctx context.Context, taskChan <-chan string, resultChan chan<- map[string]uint64) {
	res := map[string]uint64{}

	for {
		select {
		case <-ctx.Done():
			return
		case url, ok := <-taskChan:
			if !ok {
				return
			}

			essay, err := sources.LoadEngadgetBlogPost(ctx, e.httpClient, url)
			if err != nil {
				e.logger.Error("Failed to parse essay", "url", url)
				continue
			}

			counts := countValidWords(essay, e.wordBank)
			for k, v := range counts {
				res[k] += v
			}

			e.logger.Info("Successfully processed essay", "url", url)

			resultChan <- res
		}
	}
}

func streamFileToChannel(ctx context.Context, path string, out chan<- string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case out <- scanner.Text(): // blocking when buffer is full
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to iterate words: %w", err)
	}

	return nil
}

func countValidWords(essay *sources.Essay, bank *words.Bank) map[string]uint64 {
	counts := make(map[string]uint64)
	wordsList := strings.Fields(essay.Title + " " + essay.Description + " " + essay.Content)

	for _, w := range wordsList {
		cleanWord := normalizeWord(w)
		if bank.IsValid(cleanWord) {
			counts[cleanWord]++
		}
	}

	return counts
}

// lowercase and removes first and last characters if punctuation, besides apostrophe at the end
func normalizeWord(word string) string {
	if len(word) == 0 {
		return ""
	}

	first := word[0]
	if unicode.IsPunct(rune(first)) {
		word = word[1:]
	}

	last := word[len(word)-1]
	if unicode.IsPunct(rune(last)) && last != '\'' {
		word = word[:len(word)-1]
	}

	return strings.ToLower(word)
}

func mergeCounts(m1, m2 map[string]uint64) map[string]uint64 {
	result := make(map[string]uint64)

	for k, v := range m1 {
		result[k] = v
	}

	for k, v := range m2 {
		result[k] += v
	}

	return result
}

// ordered by count desc, then alphabetical order
func topNWords(counts map[string]uint64, n int) []WordCount {
	wc := make([]WordCount, 0, len(counts))
	for word, count := range counts {
		wc = append(wc, WordCount{word, count})
	}
	sort.Slice(wc, func(i, j int) bool {
		return wc[i].Count > wc[j].Count // descending
	})
	if len(wc) > n {
		wc = wc[:n]
	}
	return wc
}
