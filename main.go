package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
)

const (
	defaultWorkerCount = 10
)

func normalizeURL(url string) string {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return fmt.Sprintf("http://%s", url)
	}
	return url
}

func fetchMD5(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	hash := md5.New()
	_, err = io.Copy(hash, resp.Body)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

type task struct {
	url string
}

type result struct {
	url string
	hash string
}

type app struct {
	workerCount int
	taskCount int
	tasks chan task
	results chan result
	errors chan error
}

func worker(id int, tasks chan task, results chan result, errors chan error, wg *sync.WaitGroup) {
	for t := range tasks {
		normalizedURL := normalizeURL(t.url)
		log.Printf("Worker %d processing task with url '%s' (normalized to '%s')", id, t.url, normalizedURL)
		hash, err := fetchMD5(normalizedURL)
		if err != nil {
			errors <- err
			return
		}
		results <- result{url: normalizedURL, hash: hash}
	}
	wg.Done()
}

func (a *app) Collect() ([]result, error) {
	results := make([]result, 0, a.taskCount)
	wg := sync.WaitGroup{}
	for i := 1; i <= a.workerCount; i++ {
		wg.Add(1)
		go worker(i, a.tasks, a.results, a.errors, &wg)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		done <- struct{}{}
	}()

	for {
		select {
		case err := <- a.errors:
			return nil, err
		case result := <- a.results:
			results = append(results, result)
		case <- done:
			return results, nil
		}
	}
}

func NewApp(urls []string, workerCount int) app {
	taskCount := len(urls)
	tasks := make(chan task, taskCount)
	for i := range urls {
		tasks <- task{url: urls[i]}
	}
	close(tasks)

	// sanity checks
	if workerCount > taskCount {
		workerCount = taskCount
	}
	if workerCount <= 0 {
		workerCount = defaultWorkerCount
	}

	results := make(chan result, workerCount)
	errors := make(chan error)

	return app{workerCount: workerCount, taskCount: taskCount, tasks: tasks, results: results, errors: errors}
}

func PrettyPrint(results []result) {
	for _, result := range results {
		fmt.Printf("%s %s\n", result.url, result.hash)
	}
}

func main() {
	parallel := flag.Int("parallel", defaultWorkerCount, "number of concurrent requests")
	flag.Parse()
	app := NewApp(flag.Args(), *parallel)
	results, err := app.Collect()
	if err != nil {
		panic(err)
	}
	PrettyPrint(results)
}
