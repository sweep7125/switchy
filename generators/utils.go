// generators/utils.go
package generators

import (
	"bufio"
	"fmt"
	"net/http"
	"sort"
	"strings"
)

// FetchDomains загружает данные из raw ссылок и возвращает список доменов.
func FetchDomains(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("ошибка при загрузке URL %s: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неудачный HTTP статус для URL %s: %s", url, resp.Status)
	}

	var domains []string
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// Убираем комментарии после домена
		domain := strings.SplitN(line, "#", 2)[0]
		domain = strings.TrimSpace(domain)
		if domain != "" {
			domains = append(domains, domain)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при чтении ответа для URL %s: %v", url, err)
	}

	return domains, nil
}

// OptimizeDomains удаляет дубликаты и оптимизирует список доменов.
func OptimizeDomains(domains []string) []string {
	cleanedDomains := make(map[string]struct{})
	for _, domain := range domains {
		cleanedDomains[domain] = struct{}{}
	}

	// Преобразуем карту обратно в срез.
	uniqueDomains := make([]string, 0, len(cleanedDomains))
	for domain := range cleanedDomains {
		uniqueDomains = append(uniqueDomains, domain)
	}

	// Сортируем домены по количеству точек (от меньшего к большему).
	sort.Slice(uniqueDomains, func(i, j int) bool {
		return strings.Count(uniqueDomains[i], ".") < strings.Count(uniqueDomains[j], ".")
	})

	optimizedDomains := make(map[string]struct{})
	for _, domain := range uniqueDomains {
		covered := false
		for existing := range optimizedDomains {
			if domain == existing || strings.HasSuffix(domain, "."+existing) {
				covered = true
				break
			}
		}
		if !covered {
			optimizedDomains[domain] = struct{}{}
		}
	}

	// Преобразуем карту обратно в срез и сортируем.
	finalDomains := make([]string, 0, len(optimizedDomains))
	for domain := range optimizedDomains {
		finalDomains = append(finalDomains, domain)
	}

	sort.Strings(finalDomains)

	return finalDomains
}

// GenerateSwitchyOmegaFormat преобразует список доменов в формат SwitchyOmega.
func GenerateSwitchyOmegaFormat(domains []string) string {
	var builder strings.Builder
	builder.WriteString("#BEGIN\n\n[Wildcard]\n")
	for _, domain := range domains {
		builder.WriteString(fmt.Sprintf("*://*.%s/*\n", domain))
	}
	builder.WriteString("#END")
	return builder.String()
}
