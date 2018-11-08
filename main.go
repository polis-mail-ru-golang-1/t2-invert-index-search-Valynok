package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/t2-invert-index-search-Valynok/handlers"
	zap "go.uber.org/zap"
)

func main() {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{
		"./myproject.log",
	}
	logger, _ := cfg.Build()
	handlers.Logger = logger.Sugar()
	filesDirectoryPath := os.Args[1]
	handlers.IndexDirectory(filesDirectoryPath)
	http.HandleFunc("/", handlers.IndexHandler)
	http.HandleFunc("/search", handlers.SearchHandler)
	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}
