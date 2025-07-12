package kafka

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/tousart/mailer/repository"
)

type Worker struct {
	recipient repository.Recipient
}

func NewWorker(recipient repository.Recipient) *Worker {
	return &Worker{recipient: recipient}
}

func (w *Worker) Mail(ctx context.Context, wg *sync.WaitGroup, errorChan chan error) {
	defer wg.Done()
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				key, value, err := w.recipient.ReceiveMessage(ctx)
				if err != nil {
					log.Printf("receive error: %v\n", err)
					errorChan <- err
					return
				}

				err = SendMessage(key, value)
				if err != nil {
					log.Printf("sending error: %v\n", err)
					errorChan <- err
					return
				}
			}
		}
	}()
}

func SendMessage(key, value string) error {
	go func() {
		fmt.Printf("Отправка письма... (key: %s, value: %s)\n", key, value)
		time.Sleep(time.Second * 10)
		fmt.Println("Письмо отправлено")
	}()
	return nil
}
