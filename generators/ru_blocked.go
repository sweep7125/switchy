// generators/ru_blocked.go
package generators

import (
	"fmt"
)

// GenerateRuBlocked генерирует профиль ru-blocked.
func GenerateRuBlocked() (string, error) {
	// Список ссылок на файлы с доменами.
	fileURLs := []string{
		"https://community.antifilter.download/list/domains.lst",
		"https://raw.githubusercontent.com/1andrevich/Re-filter-lists/refs/heads/main/domains_all.lst",
		// Добавляйте сюда новые URL по необходимости
	}

	// Сбор всех доменов из файлов.
	allDomains, err := FetchAllDomains(fileURLs)
	if err != nil {
		return "", fmt.Errorf("ошибка при сборе доменов: %v", err)
	}

	// Оптимизация доменов.
	optimizedDomains := OptimizeDomains(allDomains)

	// Преобразование доменов в формат SwitchyOmega.
	switchyOmegaFormat := GenerateSwitchyOmegaFormat(optimizedDomains)

	return switchyOmegaFormat, nil
}

// FetchAllDomains загружает домены из всех указанных URL.
func FetchAllDomains(urls []string) ([]string, error) {
	var allDomains []string
	for _, url := range urls {
		domains, err := FetchDomains(url)
		if err != nil {
			return nil, err
		}
		allDomains = append(allDomains, domains...)
	}
	return allDomains, nil
}
