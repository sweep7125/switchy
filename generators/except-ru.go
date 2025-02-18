// generators/except_ru.go
package generators

import (
	"fmt"
	"strings"
)

// GenerateExceptRu генерирует профиль except-ru.
func GenerateExceptRu() (string, error) {
	// Исходная ссылка на файл.
	baseURL := "https://raw.githubusercontent.com/v2fly/domain-list-community/refs/heads/master/data/category-ru"

	// Рекурсивная обработка доменов из файла и всех директив include.
	domainsMap, err := ProcessFile(baseURL, make(map[string]struct{}))
	if err != nil {
		return "", fmt.Errorf("ошибка при обработке файла %s: %v", baseURL, err)
	}

	// Преобразовать map в slice.
	var allDomains []string
	for domain := range domainsMap {
		allDomains = append(allDomains, domain)
	}

	// Оптимизация доменов.
	optimizedDomains, err := OptimizeDomains(allDomains)
	if err != nil {
		return "", fmt.Errorf("ошибка оптимизации доменов: %v", err)
	}

	// Преобразование доменов в формат SwitchyOmega.
	switchyOmegaFormat := GenerateSwitchyOmegaFormat(optimizedDomains)

	return switchyOmegaFormat, nil
}

// ProcessFile рекурсивно обрабатывает файл и все включенные файлы.
func ProcessFile(url string, processedURLs map[string]struct{}) (map[string]struct{}, error) {
	if _, exists := processedURLs[url]; exists {
		// Файл уже обработан.
		return make(map[string]struct{}), nil
	}

	processedURLs[url] = struct{}{}

	domains, err := FetchFileContent(url)
	if err != nil {
		return nil, err
	}

	result := make(map[string]struct{})
	for _, line := range domains {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Убираем комментарии из строки.
		if idx := strings.Index(line, "#"); idx != -1 {
			line = strings.TrimSpace(line[:idx])
		}

		// Проверяем директиву include.
		if strings.HasPrefix(line, "include:") {
			includeFile := strings.TrimSpace(strings.TrimPrefix(line, "include:"))
			includeURL := fmt.Sprintf("https://raw.githubusercontent.com/v2fly/domain-list-community/refs/heads/master/data/%s", includeFile)
			includedDomains, err := ProcessFile(includeURL, processedURLs)
			if err != nil {
				return nil, err
			}
			for domain := range includedDomains {
				result[domain] = struct{}{}
			}
		} else {
			result[line] = struct{}{}
		}
	}

	return result, nil
}

// FetchFileContent загружает данные из raw ссылок и возвращает список доменов.
func FetchFileContent(url string) ([]string, error) {
	return FetchDomains(url)
}
