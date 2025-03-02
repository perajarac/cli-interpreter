package main

import (
	"os"
	"os/signal"
	"syscall"

	read "github.com/perajarac/cli-interpreter/reader"
)

var reader *read.Reader = read.NewReader()

func main() {
	read.SetUpUser()

	sigChan := make(chan os.Signal, 1)
	// Notify on Interrupt (Ctrl+C) and SIGTERM.
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		cleanupDone := make(chan struct{})

		// Listen for signals.
		go func() {
			<-sigChan
			reader.Clear()

			close(cleanupDone)
		}()

		// Wait for cleanup to finish before exiting.
		<-cleanupDone
		os.Exit(0)
	}()

	for {
		reader.MainLoop()
	}

}
