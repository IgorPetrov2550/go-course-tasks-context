// Задание 5: Дочерний контекст не отменяет родителя
//
// Важная деталь: контексты образуют дерево.
//   - Отмена родителя отменяет всех детей
//   - Отмена ребёнка НЕ отменяет родителя
//
// Это позволяет, например, дать отдельному запросу свой таймаут,
// не затрагивая другие запросы и глобальное состояние.
//
// Схема:
//
//   context.Background()          ← корень, никогда не отменяется
//       └── appCtx (WithCancel)   ← "жизнь приложения"
//               ├── req1Ctx (WithTimeout 500ms)  ← таймаут запроса 1
//               └── req2Ctx (WithTimeout 2s)     ← таймаут запроса 2
//
// Когда req1Ctx истекает - req2Ctx и appCtx НЕ отменяются.
// Когда appCtx отменяется - все дочерние контексты тоже отменяются.
//
// ──────────────────────────────────────────────────────────────────────────────
//
// Напиши функцию handleRequest(ctx context.Context, name string, delay time.Duration),
// которая:
//   - создаёт дочерний контекст с таймаутом из параметра delay:
//       reqCtx, cancel := context.WithTimeout(ctx, delay)
//       defer cancel()
//   - имитирует долгую работу (1 секунда) через select + time.After + reqCtx.Done()
//   - если успела - выводит "<name>: запрос выполнен"
//   - если таймаут - выводит "<name>: таймаут запроса (%v)", reqCtx.Err()
//
// В main():
//   1. Создай appCtx, appCancel := context.WithCancel(context.Background())
//   2. Создай WaitGroup для двух горутин
//   3. Запусти горутину: handleRequest(appCtx, "запрос-1", 500*time.Millisecond)
//      - таймаут 500мс, работа 1с → не успеет
//   4. Запусти горутину: handleRequest(appCtx, "запрос-2", 2*time.Second)
//      - таймаут 2с, работа 1с → успеет
//   5. wg.Wait()
//   6. Выведи "оба запроса завершены, приложение продолжает работу"
//   7. Вызови appCancel() - это отменит appCtx
//   8. Выведи "приложение завершено"
//
// Ожидаемый вывод:
//   запрос-1: таймаут запроса (context deadline exceeded)
//   запрос-2: запрос выполнен
//   оба запроса завершены, приложение продолжает работу
//   приложение завершено
//
// Обрати внимание: таймаут запрос-1 НЕ влияет на запрос-2 и на appCtx.
//
// Запусти: go run main.go

package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// TODO: напиши функцию handleRequest(ctx context.Context, name string, delay time.Duration)

func handleRequest(ctx context.Context, name string, delay time.Duration) {
	reqCtx, cansel := context.WithTimeout(ctx, delay)
	defer cansel()
	select {
	case <-time.After(1 * time.Second):
		fmt.Printf("%s : запрос выполнен\n", name)
	case <-reqCtx.Done():
		fmt.Printf("%s : таймаут запроса (%v)\n", name, reqCtx.Err())
	}
}

func main() {
	// TODO: создай appCtx с отменой
	// appCtx, appCancel := context.WithCancel(context.Background())
	appCtx, appCancel := context.WithCancel(context.Background())
	// TODO: создай WaitGroup локально
	var wg sync.WaitGroup

	// TODO: запусти два запроса в горутинах
	// go func() { wg.Add(1) ... handleRequest(appCtx, "запрос-1", 500*time.Millisecond) }()
	// go func() { wg.Add(1) ... handleRequest(appCtx, "запрос-2", 2*time.Second) }()
	wg.Add(1)
	go func() {
		defer wg.Done()
		handleRequest(appCtx, "запрос-1", 500*time.Millisecond)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		handleRequest(appCtx, "запрос-2", 2*time.Second)
	}()
	// TODO: wg.Wait()
	wg.Wait()
	// TODO: выведи "оба запроса завершены, приложение продолжает работу"
	fmt.Println("оба запроса завершены, приложение продолжает работу")
	// TODO: appCancel()
	appCancel()
	// TODO: выведи "приложение завершено"
	fmt.Println("приложение завершено")

}
