package main

import (
	"avito2024/internal/app"
	_ "github.com/lib/pq"
	"os"
)

func main() {
	r := app.SetupRouter()
	r.Run(os.Getenv("PORT"))
}
