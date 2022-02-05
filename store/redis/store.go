package redis

import (
	"alibaba/shortener/algo"
	"alibaba/shortener/domain"
	"context"
	"fmt"
	"github.com/fatih/structs"
	rds "github.com/go-redis/redis"
	"math/rand"
)

type redisStore struct {
	client *rds.Client
}

func NewRedisStore() *redisStore {
	c := rds.NewClient(&rds.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	// check ping
	if err := c.Ping().Err(); err != nil {
		panic(err)
	}
	return &redisStore{
		client: c,
	}
}

func (r *redisStore) GetOriginal(ctx context.Context, short string) (original string, err error) {
	randNum, err := algo.Decode(short)
	if err != nil {
		return "", domain.ErrBadParamInput
	}
	defer func() {
		if err == nil {
			r.client.HIncrBy(fmt.Sprintf("short:%d", randNum), "Views", 1)
		}
	}()
	original, err = r.client.WithContext(ctx).HGet(fmt.Sprintf("short:%d", randNum), "Original").Result()
	if err == rds.Nil {
		return "", domain.ErrNotFound
	} else if err != nil {
		return "", domain.ErrInternalServerError
	}
	return original, nil
}

func (r *redisStore) GetInfo(ctx context.Context, short string) (info map[string]string, err error) {
	randNum, err := algo.Decode(short)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	if !r.isUsed(ctx, randNum) {
		return nil, domain.ErrNotFound
	}
	return r.client.WithContext(ctx).HGetAll(fmt.Sprintf("short:%d", randNum)).Result()
}

func (r *redisStore) Save(ctx context.Context, original string) (short string, err error) {
	var randNum uint64
	for used := true; used; used = r.isUsed(ctx, randNum) {
		randNum = rand.Uint64() + 1e6
	}
	short = algo.Encode(randNum)
	record := domain.Record{
		Original:  original,
		Shortened: short,
		Views:     0,
	}
	err = r.client.WithContext(ctx).HMSet(fmt.Sprintf("short:%d", randNum), structs.Map(record)).Err()
	if err != nil {
		return "", domain.ErrInternalServerError
	}
	return "http://127.0.0.1:8080/sh/" + algo.Encode(randNum), nil
}

func (r *redisStore) isUsed(ctx context.Context, num uint64) (isUsed bool) {
	exists, err := r.client.WithContext(ctx).Exists(fmt.Sprintf("short:%d", num)).Result()
	if err != nil {
		return false
	}
	return exists > 0
}
