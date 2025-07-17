package kafka

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"sync"

	"github.com/tousart/mailer/domain/models"
	"github.com/tousart/mailer/repository"
	"gopkg.in/gomail.v2"
)

type Worker struct {
	recipient repository.Recipient
}

func NewWorker(recipient repository.Recipient) *Worker {
	return &Worker{recipient: recipient}
}

func (w *Worker) Mail(ctx context.Context, errorChan chan error, wg *sync.WaitGroup) {
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				log.Println("context has been closed")
				return
			case err := <-errorChan:
				log.Printf("sending error: %v", err)
				return
			default:
				key, value, err := w.recipient.ReceiveMessage(context.Background())
				if err != nil {
					log.Printf("receive error: %v\n", err)
					errorChan <- err
					return
				}

				wg.Add(1)
				go sendMessage(key, value, errorChan, wg)
			}
		}
	}()
}

func sendMessage(key string, value []byte, errorChan chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Printf("sending message to user... (key: %s)\n", key)

	var email models.Email
	err := json.Unmarshal(value, &email)
	if err != nil {
		errorChan <- err
		return
	}

	from := os.Getenv("POST_NAME")
	password := os.Getenv("POST_PASSWORD")
	to := email.Email

	subject := "Топ 1 сервис на Марсе"
	body := "Привет, " + email.Login + "! Рад видеть тебя в этом притоне)\nЗаходи к нам почаще, обнимаю."

	message := gomail.NewMessage()
	message.SetHeader("From", from)
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)
	message.SetBody("text/plain", body)

	d := gomail.NewDialer("smtp.yandex.ru", 587, from, password)

	if err := d.DialAndSend(message); err != nil {
		errorChan <- err
		return
	}

	log.Println("message has been sent to user")
}
