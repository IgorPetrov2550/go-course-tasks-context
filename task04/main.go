// Задание 4: Контекст + горутины + WaitGroup вместе
//
// Реальная задача: запустить несколько горутин, каждая делает свою работу,
// и всех остановить по таймауту - аккуратно.
//
// Здесь важно понять как контекст, WaitGroup и горутины работают ВМЕСТЕ:
//
//   context.WithTimeout - даёт дедлайн (когда всё должно завершиться)
//   WaitGroup           - позволяет дождаться ВСЕХ горутин
//   горутины            - делают параллельную работу
//
// Обрати внимание: WaitGroup создаётся в main() и передаётся как параметр.
// Каждая функция которая запускает горутину - должна получить wg явно.
// Это та же идея что в задании 3 про замыкания.
//
// ──────────────────────────────────────────────────────────────────────────────
//
// Есть список "сервисов" которые нужно проверить параллельно:
//   services := []string{"users", "orders", "payments", "inventory", "notifications"}
//
// Напиши функцию checkService(ctx context.Context, name string, wg *sync.WaitGroup),
// которая:
//   - вызывает defer wg.Done()
//   - имитирует проверку: time.Sleep(случайное время от 100 до 600мс)
//     задержка: time.Duration(100 + rand.Intn(500)) * time.Millisecond
//   - если контекст отменился ДО завершения - выводит
//     "<name>: проверка отменена (таймаут)"
//   - иначе выводит "<name>: ОК"
//   - используй select для проверки завершения
//
// В main():
//   1. Создай контекст с таймаутом 300мс (часть сервисов успеет, часть нет)
//   2. defer cancel()
//   3. Создай WaitGroup локально
//   4. Запусти checkService для каждого сервиса в горутине
//   5. wg.Wait() - дождись всех
//   6. Выведи "проверка завершена"
//
// Ожидаемый вывод (порядок и результаты меняются каждый раз):
//   users: ОК
//   orders: проверка отменена (таймаут)
//   payments: ОК
//   inventory: проверка отменена (таймаут)
//   notifications: проверка отменена (таймаут)
//   проверка завершена
//
// Запусти несколько раз - результаты будут разными.
//
// Запусти: go run main.go

package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// TODO: напиши функцию checkService(ctx context.Context, name string, wg *sync.WaitGroup)
// Подсказка:
//   delay := time.Duration(100+rand.Intn(500)) * time.Millisecond
//   select {
//   case <-time.After(delay):
//       // успех
//   case <-ctx.Done():
//       // таймаут
//   }

func checkService(ctx context.Context, name string, wg *sync.WaitGroup) {
	defer wg.Done()
	delay := time.Duration(100+rand.Intn(500)) * time.Millisecond
	select {
	case <-time.After(delay):
		fmt.Printf("%s: OK\n", name) // успех
	case <-ctx.Done():
		fmt.Printf("%s: проверка отменена (таймаут)\n", name) // таймаут
	}
}

func main() {
	services := []string{"users", "orders", "payments", "inventory", "notifications"}

	// TODO: создай контекст с таймаутом 300мс
	ctx, cansel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cansel()
	// TODO: создай WaitGroup локально
	var wg sync.WaitGroup
	// TODO: запусти checkService для каждого сервиса
	for _, service := range services {
		wg.Add(1)
		go checkService(ctx, service, &wg)
	}
	// TODO: wg.Wait() и выведи "проверка завершена"
	wg.Wait()
	fmt.Println("проверка завершена")
}
