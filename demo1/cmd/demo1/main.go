package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/heetch/go-meetup-plugins-2018/demo1"
	"github.com/pkg/errors"
)

func main() {
	err := start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
	}
}

func start() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	fmt.Println("Core: Loading plugins...")
	loader := demo1.PluginLoader{Path: "./bin", PlayerAddr: "some-addr"}
	err := loader.Load()
	if err != nil {
		return errors.Wrap(err, "player: failed to load plugins")
	}
	fmt.Println("Core: All plugins successfully loaded")

	<-quit
	fmt.Println()
	fmt.Println("Core: closing plugins")
	return loader.Stop()
}
