package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

func main() {
	const maxTasks = 5
	const searchText = "Go"

	var urls = []string{
		"https://golang.org/",
		"https://ru.wikipedia.org/wiki/Go",
		"https://blog.skillfactory.ru/glossary/golang/",
		"https://hub.docker.com/_/golang",
		"https://go.dev/",
		"https://github.com/golang/go",
		"https://golangify.com/",
	}

	var wg = sync.WaitGroup{}
	var urlsChannel = make(chan string)
	var resultsChannel = make(chan UrlCountResult, len(urls))
	var minTasks = min(maxTasks, len(urls))

	for i := 0; i < minTasks; i++ {
		wg.Add(1)
		go handleTask(searchText, &wg, urlsChannel, resultsChannel)
	}

	for _, url := range urls {
		urlsChannel <- url
	}

	close(urlsChannel)

	var total = 0
	for i := 0; i < len(urls); i++ {
		var result = <-resultsChannel
		total += result.count
		fmt.Printf("Count for %s: %d\n", result.url, result.count)
	}

	wg.Wait()

	fmt.Printf("Total: %d", total)
}

func handleTask(searchText string, wg *sync.WaitGroup, urls <-chan string, results chan<- UrlCountResult) {
	for url := range urls {
		queryCountOfText(url, searchText, results)
	}
	wg.Done()
}

func queryCountOfText(url string, searchText string, results chan<- UrlCountResult) {
	var httpClient = http.Client{}
	var count = 0

	var getResult = func() UrlCountResult {
		return UrlCountResult{url, count}
	}

	response, err := httpClient.Get(url)
	if err != nil {
		println(err.Error())
		results <- getResult()
		return
	}

	bodyText, err := io.ReadAll(response.Body)
	if err != nil {
		println(err.Error())
		results <- getResult()
		return
	}

	count = strings.Count(string(bodyText), searchText)

	results <- getResult()
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

type UrlCountResult struct {
	url   string
	count int
}
