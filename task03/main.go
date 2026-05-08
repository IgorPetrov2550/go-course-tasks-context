// Задание 3: Передача данных через контекст
//
// context.WithValue позволяет "положить" значение в контекст и достать его
// в любой функции глубже по стеку вызовов - без явной передачи параметром.
//
// Типичный случай: request ID, user ID, trace ID.
// Их нужно везде логировать, но добавлять в каждую функцию как параметр - громоздко.
//
// ПРАВИЛА использования context.WithValue:
//
//   1. Ключ должен быть СВОИМ типом - не строкой, не int.
//      Почему? Чтобы избежать коллизий с другими пакетами:
//
//        // ПЛОХО:
//        ctx = context.WithValue(ctx, "user_id", 42)
//        // Другой пакет тоже может использовать ключ "user_id"!
//
//        // ХОРОШО:
//        type contextKey string
//        const userIDKey contextKey = "user_id"
//        ctx = context.WithValue(ctx, userIDKey, 42)
//        // Теперь ключ уникален для этого пакета
//
//   2. context.WithValue - только для данных которые нужны "везде" (request-scoped).
//      Бизнес-параметры (количество, сумма, имя) передавай обычными аргументами.
//
// ──────────────────────────────────────────────────────────────────────────────
//
// Реализуй цепочку: main → handler → service → repository.
//
// Объяви:
//   type contextKey string
//   const requestIDKey contextKey = "request-id"
//   const userIDKey    contextKey = "user-id"
//
// В handler(ctx context.Context, userID int):
//   - добавь в контекст requestID = fmt.Sprintf("req-%d", userID)
//   - добавь в контекст userID
//   - вызови service(ctx) и выведи ошибку если есть
//
// В service(ctx context.Context) error:
//   - достань requestID из контекста
//   - выведи "[<requestID>] service: обрабатываем пользователя"
//   - вызови repository(ctx) и оберни ошибку: fmt.Errorf("service: %w", err)
//
// В repository(ctx context.Context) error:
//   - достань requestID и userID из контекста
//   - выведи "[<requestID>] repo: запрос для user=<userID>"
//   - если userID < 0 - верни errors.New("user not found")
//   - иначе выведи "[<requestID>] repo: готово" и верни nil
//
// В main():
//   - вызови handler(context.Background(), 42)  - успешно
//   - вызови handler(context.Background(), -1)  - ошибка
//
// Ожидаемый вывод:
//   [req-42] service: обрабатываем пользователя
//   [req-42] repo: запрос для user=42
//   [req-42] repo: готово
//   [req--1] service: обрабатываем пользователя
//   [req--1] repo: запрос для user=-1
//   ошибка в handler: service: user not found
//
// Запусти: go run main.go

package main

import (
	"context"
	"errors"
	"fmt"
)

type contextKey string

const requestIDKey contextKey = "request-id"
const userIDKey contextKey = "user-id"

// TODO: реализуй repository(ctx context.Context) error

func repository(ctx context.Context) error {
	reqID := ctx.Value(requestIDKey).(string)
	userID := ctx.Value(userIDKey).(int)

	fmt.Printf("[%s] repo: запрос для user=%d\n", reqID, userID)

	if userID < 0 {
		return errors.New("user not found")
	}

	fmt.Printf("[%s] repo: готово\n", reqID)
	return nil
}

// TODO: реализуй service(ctx context.Context) error

func service(ctx context.Context) error {
	reqID := ctx.Value(requestIDKey).(string)

	fmt.Printf("[%s] service: обрабатываем пользователя\n", reqID)

	err := repository(ctx)
	if err != nil {
		return fmt.Errorf("service: %w", err)
	}

	return nil
}

// TODO: реализуй handler(ctx context.Context, userID int)

func handler(ctx context.Context, userID int) {
	reqID := fmt.Sprintf("req-%d", userID)
	ctx = context.WithValue(ctx, requestIDKey, reqID)
	ctx = context.WithValue(ctx, userIDKey, userID)

	err := service(ctx)
	if err != nil {
		fmt.Printf("ошибка в handler: %v\n", err)
	}
}

func main() {
	handler(context.Background(), 42)
	handler(context.Background(), -1)
}
