package kafka

import (
	"context"
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
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				key, value, err := w.recipient.ReceiveMessage(context.Background())
				if err != nil {
					log.Printf("receive error: %v\n", err)
					errorChan <- err
					return
				}

				wg.Add(1)
				err = SendMessage(key, value, wg)
				if err != nil {
					log.Printf("sending error: %v\n", err)
					errorChan <- err
					return
				}
			}
		}
	}()
}

func SendMessage(key, value string, wg *sync.WaitGroup) error {
	go func() {
		defer wg.Done()
		log.Printf("Отправка письма... (key: %s, value: %s)\n", key, value)
		time.Sleep(time.Second * 10)
		log.Println("Письмо отправлено")
	}()
	return nil
}
