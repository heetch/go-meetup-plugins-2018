package main

import (
	"fmt"
	"os"
	"os/signal"
)

func main() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	fmt.Println("Plugin2: loaded")
	<-quit
	fmt.Println("Plugin2: closing")
}
