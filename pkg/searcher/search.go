package searcher

import (
	"io/fs"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"word-search-in-files/pkg/internal/dir"
)

type Searcher struct {
	FS fs.FS
}

func (s *Searcher) Search(word string) (files []string, err error) {

	// Слайс списка всех файлов в текущей директории
	allFiles, err := dir.FilesFS(s.FS, ".")
	if err != nil {
		return nil, err
	}

	// Регулярное выражение для точного поиска слова, можно было бы использовать strings.Contains,
	// но тогда выходили бы все слова, которые так же включают частичное сопадение искомого слова
	// так же не работает с кирилицей такое regexPattern := `\b` + regexp.QuoteMeta(word) + `\b`
	// потому добавила условия слева начало строки или пробельный символ,
	// справа конец строки, или пробельный символ, или точка, или запятая
	regexPattern := `(?:^|\s)` + regexp.QuoteMeta(word) + `(?:\s|$|,|.)`
	re, err := regexp.Compile(regexPattern)
	if err != nil {
		return nil, err
	}

	results := make(chan string) // Канал для передачи результатов
	errs := make(chan error)     // Канал для передачи ошибок
	var wg sync.WaitGroup        // Для ожидания завершения всех горутин

	// Запуск поиска в каждом файле с помощью горутин
	for _, path := range allFiles {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			// Читаем содержимое файла
			content, err := fs.ReadFile(s.FS, path)
			if err != nil {
				errs <- err
				return
			}

			// Проверка содержимого файла
			if re.Match(content) {
				// Убираем из имени файла расширение
				filename := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
				results <- filename
			}
		}(path)
	}

	// Ожидаем завершение всех горутин и закрываем каналы
	go func() {
		wg.Wait()
		close(results)
		close(errs)
	}()

	// Сбор результатов
	for result := range results {
		files = append(files, result)
	}
	sort.Strings(files) //сортируем список файлов

	// Проверка ошибок
	for err := range errs {
		if err != nil {
			return nil, err
		}
	}

	return files, nil
}
