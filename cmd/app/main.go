package main

import (
	"RomManager/internal/app"
	"RomManager/internal/config"
	"log"
	_ "net/http/pprof"
)

func main() {
	configPath := "config.yml"
	c, err := config.New(configPath)
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}
	log.Printf("loaded config: %s", configPath)

	a, err := app.New(c)
	if err != nil {
		log.Fatalf("could not create app: %v", err)
	}
	defer a.Destroy()

	a.Run()
}
