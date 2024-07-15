package main

import (
	"fmt"
	"sso/internal/config"
)

// собирает в себе все модули
func main() {
	cfg := config.MustLoad()

	fmt.Println(cfg)

	// TODO: Инициализировать логгер

	// TODO: Инициализировать  приложение (/app)

	// TODO: запустить gRPC-сервер приложения
}
