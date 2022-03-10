package main

import (
	"math"
	"sync"
	"testing"
)

const (
	exampleComHash = "84238dfc8092e5d9c0dac8ef93371a07"
)

var sitesToTest = []string{"adjust.com", "google.com", "facebook.com", "yahoo.com", "yandex.com", "twitter.com",
	"reddit.com/r/funny", "reddit.com/r/notfunny", "baroquemusiclibrary.com", "http://example.com"}
var workerCounts = []int{1, 3, 200, math.MaxInt64, 20, 10, len(sitesToTest), -1, 0}

func failOnErr(err error, t *testing.T) {
	if err != nil {
		t.Fatalf("Error while collecting results: %v", err)
	}
}

// unreliable test, example.com might change in the future. but it works fine in the scope of the test task
func TestKnownHash(t *testing.T) {
	a := NewApp([]string{"example.com"}, 1)
	results, err := a.Collect()
	failOnErr(err, t)
	if results[0].hash != exampleComHash {
		t.Fatalf("Unexpected hash for %s. Expected %s, got %s", results[0].url, exampleComHash, results[0].hash)
	}
}

func testWorker(wg *sync.WaitGroup, workerCount int, t *testing.T) {
	a := NewApp(sitesToTest, workerCount)
	_, err := a.Collect()
	failOnErr(err, t)
	wg.Done()
}

// The only thing we can reliably test is the program's ability to handle different worker counts
func TestWorkerCount(t *testing.T) {
	wg := sync.WaitGroup{}
	for i := range workerCounts {
		wg.Add(1)
		go testWorker(&wg, workerCounts[i], t)
	}
	wg.Wait()
}