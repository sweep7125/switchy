// main.go
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"profilegen/generators"
)

func main() {
	fmt.Println("Начинаем процесс генерации профилей SwitchyOmega...")

	// Создание папки results, если она не существует
	resultsDir := "results"
	if err := os.MkdirAll(resultsDir, os.ModePerm); err != nil {
		log.Fatalf("Не удалось создать папку %s: %v", resultsDir, err)
	}

	// Список генераторов
	generatorsList := []struct {
		Name         string
		GenerateFunc func() (string, error)
	}{
		{
			Name:         "ru-blocked",
			GenerateFunc: generators.GenerateRuBlocked,
		},
		{
			Name:         "except-ru",
			GenerateFunc: generators.GenerateExceptRu,
		},
		// Добавляйте сюда новые генераторы по мере необходимости
		// {
		//     Name:         "us-unblocked",
		//     GenerateFunc: generators.GenerateUsUnblocked,
		// },
	}

	// Генерация профилей
	for _, gen := range generatorsList {
		fmt.Printf("Генерируем профиль: %s\n", gen.Name)
		content, err := gen.GenerateFunc()
		if err != nil {
			log.Printf("Ошибка при генерации профиля %s: %v\n", gen.Name, err)
			continue
		}

		// Определение пути для сохранения файла
		outputPath := filepath.Join(resultsDir, fmt.Sprintf("%s.txt", gen.Name))

		// Сохранение контента в файл
		if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
			log.Printf("Ошибка при сохранении профиля %s: %v\n", gen.Name, err)
			continue
		}

		fmt.Printf("Профиль %s успешно сохранен в %s\n", gen.Name, outputPath)
	}

	fmt.Println("Генерация профилей завершена.")
}
