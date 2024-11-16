package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/ophum/gslb-poc/internal"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Servers []*internal.Server
}

var config Config
var store = internal.NewServerStore()

func init() {
	f, err := os.Open("config.yaml")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := yaml.NewDecoder(f).Decode(&config); err != nil {
		panic(err)
	}

	for _, s := range config.Servers {
		store.Set(s)
	}
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	shutdown, err := internal.Start(store)
	if err != nil {
		panic(err)
	}

	<-ctx.Done()

	ctx, stop = signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	if err := shutdown(ctx); err != nil {
		panic(err)
	}
}
