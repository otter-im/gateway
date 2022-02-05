package app

import (
	"context"
	"fmt"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-redis/cache/v8"
	"sync"
	"time"
)

type OtterClientStore struct {
	sync.RWMutex
}

func (s *OtterClientStore) GetByID(ctx context.Context, id string) (oauth2.ClientInfo, error) {
	s.RLock()
	defer s.RUnlock()

	token := new(Client)
	if err := RedisCache().Once(&cache.Item{
		Ctx:   ctx,
		Key:   fmt.Sprintf("client:%v", id),
		Value: token,
		TTL:   15 * time.Minute,
		Do: func(item *cache.Item) (interface{}, error) {
			return s.selectByID(ctx, id)
		},
	}); err != nil {
		return nil, err
	}
	return token, nil
}

func (s *OtterClientStore) selectByID(ctx context.Context, id string) (*Client, error) {
	client := new(Client)
	if err := Postgres().
		ModelContext(ctx, client).
		Where("id = ?", id).
		Select(); err != nil {
		return nil, err
	}
	return client, nil
}
