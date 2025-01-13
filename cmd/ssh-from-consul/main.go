package main

import (
	"log"
	"os"

	"github.com/rallvesh/ssh-from-consul/internal/app"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s [profile] <command> [args]", os.Args[0])
	}

	// Проверка на флаг --version
	if os.Args[1] == "--version" {
		app.ShowVersion()
		return
	}

	profile := "default"
	commandIndex := 1

	// Если первый аргумент - это профиль, используем его
	if len(os.Args) > 2 {
		profile = os.Args[1]
		commandIndex = 2
	}

	command := os.Args[commandIndex]

	// Обрабатываем команду
	app.HandleCommand(command, profile)
}
