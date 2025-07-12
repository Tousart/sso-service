package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	rep "github.com/tousart/mailer/repository/kafka"
	serv "github.com/tousart/mailer/server/kafka"
)

func main() {
	// signal for Graceful Shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	errorChan := make(chan error, 1)

	recipient := rep.NewKafkaRecipient([]string{"kafka:9092"}, "email_messages", "mailers")

	wg := new(sync.WaitGroup)

	worker := serv.NewWorker(recipient)
	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)
	worker.Mail(ctx, wg, errorChan)

	defer func() {
		cancel()
		wg.Wait()
		recipient.Reader.Close()
	}()

	select {
	case <-signalChan:
		log.Print("graceful shutdown\n")
	case err := <-errorChan:
		log.Printf("worker error: %v", err)
	}
}
