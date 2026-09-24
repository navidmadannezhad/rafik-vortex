package main

import (
	"airun/code-reviewer/router"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	router.RunRouter()
}
