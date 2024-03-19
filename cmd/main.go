package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ositlar/floodcontrol/pkg/floodcontrol"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	fc, err := floodcontrol.NewFloodController("floodcontrol", "vk", "mongodb://localhost:27017", 1, time.Millisecond*300)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 10; i++ {
		result, err := fc.Check(ctx, 2)
		if err != nil {
			log.Fatal(err)
		}
		if result {
			fmt.Println("Flood control пройден")
		} else {
			fmt.Println("Flood control не пройден")
		}
		time.Sleep(time.Millisecond * 350)
	}
	cancel()
	if err := fc.Disconnect(ctx); err != nil {
		log.Fatal(err)
	}
}

// FloodControl интерфейс, который нужно реализовать.
// Рекомендуем создать директорию-пакет, в которой будет находиться реализация.
type FloodControl interface {
	// Check возвращает false если достигнут лимит максимально разрешенного
	// кол-ва запросов согласно заданным правилам флуд контроля.
	Check(ctx context.Context, userID int64) (bool, error)
}
