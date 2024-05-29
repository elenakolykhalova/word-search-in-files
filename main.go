package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"word-search-in-files/pkg/searcher"
)

// Обработчик HTTP-запросов для поиска файлов
func searchHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем параметр "word" из запроса
	word := r.URL.Query().Get("word")
	if word == "" {
		http.Error(w, "Missing search word", http.StatusBadRequest)
		return
	}

	// Создаем файловую систему из текущего каталога
	currentFS := os.DirFS(".")

	// Создаем экземпляр Searcher
	s := searcher.Searcher{FS: currentFS}

	// Проходим по файлам и ищем нужное слово
	files, err := s.Search(word)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error searching files: %v", err), http.StatusInternalServerError)
		return
	}

	// Сериализуем результат в JSON
	response, err := json.Marshal(files)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error serializing response: %v", err), http.StatusInternalServerError)
		return
	}

	// Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func main() {
	http.HandleFunc("/files/search", searchHandler)

	fmt.Println("Starting server on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
