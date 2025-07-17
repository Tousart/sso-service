package repository

import "context"

type Recipient interface {
	ReceiveMessage(ctx context.Context) (key string, value []byte, err error)
}
