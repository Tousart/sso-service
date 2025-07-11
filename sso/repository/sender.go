package repository

import "context"

type Sender interface {
	SendMessage(ctx context.Context, key, value []byte) error
}
