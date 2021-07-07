package main

import (
	"github.com/joho/godotenv"
	"github.com/tigorlazuardi/ridit/app"
)

func main() {
	_ = godotenv.Load()
	app.Start()
}
