package main

import (
	"log"

	"github.com/kevholditch/vigilant/internal/app"
)

func main() {
	app := app.NewApp()
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
