package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	log.Print("Running...")

	stopCh := make(chan bool)
	signalCh := make(chan os.Signal)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		s := <-signalCh
		log.Print("Stopping...: ", s)
		stopCh<-true
	}()

	<-stopCh
}
