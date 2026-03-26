package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"L2.16/internal/crawler"
)

func main() {
	var urlStr string
	var depth int
	var outDir string

	flag.StringVar(&urlStr, "u", "", "URL to download")
	flag.IntVar(&depth, "n", 1, "Recursion depth")
	flag.StringVar(&outDir, "o", "./res", "Output directory")
	flag.Parse()

	if urlStr == "" {
		log.Fatal("Missing -u parameter")
	}

	c := crawler.Crawler{
		// Добавили таймаут, чтобы зависшие ресурсы отваливались по ошибке, а не вешали всю программу
		Client: &http.Client{
			Timeout: 15 * time.Second,
		},
		Visited:  make(map[string]bool),
		SaveDir:  outDir,
		MaxDepth: depth,
	}

	fmt.Printf("Start: %s (depth %d)\n", urlStr, depth)
	_, err := c.Crawl(urlStr, depth)
	if err != nil {
		log.Fatalf("Ошибка: %v\n", err)
	}
}
