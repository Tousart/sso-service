package repository

import "context"

type Recipient interface {
	ReceiveMessage(ctx context.Context) (key, value string, err error)
}
