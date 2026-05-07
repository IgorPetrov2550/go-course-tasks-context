// Задание 2: Таймаут для медленной операции
//
// context.WithTimeout - это автоматическая отмена через заданное время.
// Очень часто используется для HTTP-запросов, запросов в БД и т.д.
//
// Разница от WithCancel:
//   - WithCancel: ты вручную вызываешь cancel()
//   - WithTimeout: контекст отменится сам через N времени (и тоже можно cancel())
//
// ВАЖНО: всегда вызывай defer cancel() даже с таймаутом!
//   Если операция завершится ДО таймаута - cancel() освободит ресурсы немедленно.
//   Без defer cancel() ресурсы освободятся только когда сработает таймаут.
//
// ──────────────────────────────────────────────────────────────────────────────
//
// Напиши функцию downloadFile(ctx context.Context, url string) (string, error),
// которая имитирует скачивание файла:
//   - использует select с двумя случаями:
//       * time.After(2 * time.Second): "скачивание" заняло 2 секунды, возвращает
//         ("содержимое " + url, nil)
//       * ctx.Done(): возвращает ("", ctx.Err())
//
// В main() вызови downloadFile дважды:
//
//   Вызов 1: таймаут 3 секунды → должен успеть
//     ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
//     defer cancel()
//     result, err := downloadFile(ctx, "https://example.com/file.txt")
//     // ожидаемый вывод: скачано: содержимое https://example.com/file.txt
//
//   Вызов 2: таймаут 1 секунда → не успеет, таймаут
//     // ожидаемый вывод: ошибка: context deadline exceeded
//
// После каждого вызова выводи время выполнения через time.Since(start).
//
// Ожидаемый вывод:
//   скачано: содержимое https://example.com/file.txt (за ~2s)
//   ошибка: context deadline exceeded (за ~1s)
//
// Запусти: go run main.go

package main

import (
	"context"
	"fmt"
	"time"
)

// TODO: напиши функцию downloadFile(ctx context.Context, url string) (string, error)

func downloadFile(ctx context.Context, url string) (string, error) {
	select {
	case <-time.After(2 * time.Second):
		return "содержимое " + url, nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func main() {
	// Вызов 1: таймаут 3 секунды - должен успеть
	// TODO: создай контекст с таймаутом 3s, вызови downloadFile, выведи результат
	start1 := time.Now()
	ctx1, cancel1 := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel1()

	//res1, err1 := downloadFile(ctx1, "https://example.com/file.txt")
	if res1, err1 := downloadFile(ctx1, "https://example.com/file.txt"); err1 != nil {
		fmt.Printf("ошибка: %v (за %v)\n", err1, time.Since(start1))
	} else {
		fmt.Printf("скачано: %v (за %v)\n", res1, time.Since(start1))
	}
	// Вызов 2: таймаут 1 секунда - не успеет
	// TODO: создай контекст с таймаутом 1s, вызови downloadFile, выведи ошибку
	start2 := time.Now()

	ctx2, cancel2 := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel2()

	if res2, err2 := downloadFile(ctx2, "https://example.com/file.txt"); err2 != nil {
		fmt.Printf("ошибка: %v (за %v)\n", err2, time.Since(start2))
	} else {
		fmt.Printf("скачано: %v (за %v)\n", res2, time.Since(start2))
	}

}
