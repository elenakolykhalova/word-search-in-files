package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"word-search-in-files/pkg/searcher"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

// Хэндлер для поиска файлов
func (h *Handler) SearchHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем параметр "word" из запроса
	word := r.URL.Query().Get("word")
	if word == "" {
		http.Error(w, "Missing search word", http.StatusBadRequest)
		return
	}

	// Создаем файловую систему из каталога examples
	currentFS := os.DirFS("examples")

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
