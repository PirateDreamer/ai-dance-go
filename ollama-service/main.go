package main

import (
	"ollama-service/module/ollama"

	"github.com/PirateDreamer/going"
)

func main() {
	router := going.InitService()
	ollama.InitController()
	router.Run()
}
