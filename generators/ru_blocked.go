// generators/ru_blocked.go
package generators

import (
	"fmt"
)

// GenerateRuBlocked генерирует профиль ru-blocked в формате SwitchyOmega.
func GenerateRuBlocked() (string, error) {
	// Список URL с доменными списками.
	fileURLs := []string{
		"https://community.antifilter.download/list/domains.lst",
		"https://raw.githubusercontent.com/1andrevich/Re-filter-lists/refs/heads/main/domains_all.lst",
		// При необходимости добавляйте новые URL.
	}

	// Параллельно загружаем домены из всех URL.
	allDomains, err := FetchAllDomains(fileURLs)
	if err != nil {
		return "", fmt.Errorf("ошибка при сборе доменов: %v", err)
	}

	// Оптимизируем список доменов (удаляем дубликаты и группируем по регистрируемым доменам).
	optimizedDomains, err := OptimizeDomains(allDomains)
	if err != nil {
		return "", fmt.Errorf("ошибка при оптимизации доменов: %v", err)
	}

	// Преобразуем домены в формат SwitchyOmega.
	switchyOmegaFormat := GenerateSwitchyOmegaFormat(optimizedDomains)

	return switchyOmegaFormat, nil
}
