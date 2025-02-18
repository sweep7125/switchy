// generators/utils.go
package generators

import (
	"bufio"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/publicsuffix"
)

// HTTP-клиент с таймаутом для всех запросов.
var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

// FetchDomains загружает данные с указанного URL и возвращает список доменов.
func FetchDomains(url string) ([]string, error) {
	resp, err := httpClient.Get(url)
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
		// Пропускаем пустые строки и комментарии.
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// Убираем комментарии после домена.
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

// FetchAllDomains параллельно загружает домены по списку URL.
func FetchAllDomains(urls []string) ([]string, error) {
	var wg sync.WaitGroup
	domainsCh := make(chan []string, len(urls))
	errCh := make(chan error, len(urls))

	// Запускаем горутины для каждого URL.
	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			domains, err := FetchDomains(u)
			if err != nil {
				errCh <- err
				return
			}
			domainsCh <- domains
		}(url)
	}

	wg.Wait()
	close(domainsCh)
	close(errCh)

	// Если возникла ошибка хотя бы в одном запросе, возвращаем её.
	if len(errCh) > 0 {
		return nil, <-errCh
	}

	var allDomains []string
	for ds := range domainsCh {
		allDomains = append(allDomains, ds...)
	}

	return allDomains, nil
}

// OptimizeDomains группирует домены по регистрируемому домену (Effective TLD Plus One)
// и возвращает уникальный список, который позволяет покрыть все поддомены одним правилом.
func OptimizeDomains(domains []string) ([]string, error) {
	registrableMap := make(map[string]struct{})

	for _, domain := range domains {
		// Приводим домен к нижнему регистру для единообразия.
		domain = strings.ToLower(strings.TrimSpace(domain))
		if domain == "" {
			continue
		}

		// Получаем регистрируемый домен, например, для a.b.sinema2.top -> sinema2.top.
		regDomain, err := publicsuffix.EffectiveTLDPlusOne(domain)
		if err != nil {
			// Если не удалось получить регистрируемый домен, используем оригинал.
			regDomain = domain
		}
		registrableMap[regDomain] = struct{}{}
	}

	// Преобразуем карту в срез.
	optimized := make([]string, 0, len(registrableMap))
	for regDomain := range registrableMap {
		optimized = append(optimized, regDomain)
	}

	// Сортировка для стабильного вывода.
	sort.Strings(optimized)
	return optimized, nil
}

// GenerateSwitchyOmegaFormat преобразует список доменов в формат SwitchyOmega.
func GenerateSwitchyOmegaFormat(domains []string) string {
	var builder strings.Builder
	builder.WriteString("#BEGIN\n\n[Wildcard]\n")
	for _, domain := range domains {
		// Формат правила: *://*.домен/*
		builder.WriteString(fmt.Sprintf("*://*.%s/*\n", domain))
	}
	builder.WriteString("#END")
	return builder.String()
}
