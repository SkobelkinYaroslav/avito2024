package main

import (
	"avito2024/internal/app"
	_ "github.com/lib/pq"
)

func main() {
	r := app.SetupRouter()
	r.Run(":8080")
}
