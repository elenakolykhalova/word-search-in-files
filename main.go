package main

import (
	"fmt"
	"net/http"
	"word-search-in-files/pkg/handlers"
)

func main() {

	s := handlers.NewHandler()
	http.HandleFunc("/files/search", s.SearchHandler)

	fmt.Println("Starting server on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
