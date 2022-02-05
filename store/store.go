package store

import (
	"context"
)

type Store interface {
	GetOriginal(ctx context.Context, short string) (original string, err error)
	GetInfo(ctx context.Context, short string) (info map[string]string, err error)
	Save(ctx context.Context, original string) (short string, err error)
}
