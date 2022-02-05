package app

import (
	"context"
	"fmt"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-redis/cache/v8"
	"time"
)

type OtterTokenStore struct{}

func (s *OtterTokenStore) Create(ctx context.Context, info oauth2.TokenInfo) error {
	token := &Token{
		ClientID:    info.GetClientID(),
		UserID:      info.GetUserID(),
		RedirectURI: info.GetRedirectURI(),
		Scope:       info.GetScope(),
	}

	if code := info.GetCode(); code != "" {
		token.Code = code
		token.CodeChallenge = info.GetCodeChallenge()
		token.CodeChallengeMethod = info.GetCodeChallengeMethod().String()
		token.CodeCreateAt = info.GetCodeCreateAt()
		token.CodeExpiresAt = info.GetCodeCreateAt().Add(info.GetCodeExpiresIn())
	}

	if access := info.GetAccess(); access != "" {
		token.Access = info.GetAccess()
		token.AccessCreateAt = info.GetAccessCreateAt()
		token.AccessExpiresAt = info.GetAccessCreateAt().Add(info.GetAccessExpiresIn())
	}

	if refresh := info.GetRefresh(); refresh != "" {
		token.Refresh = info.GetRefresh()
		token.RefreshCreateAt = info.GetRefreshCreateAt()
		token.RefreshExpiresAt = info.GetRefreshCreateAt().Add(info.GetRefreshExpiresIn())
	}

	if _, err := Postgres().ModelContext(ctx, token).Insert(); err != nil {
		return err
	}
	return nil
}

func (s *OtterTokenStore) RemoveByCode(ctx context.Context, code string) error {
	if err := RedisCache().Delete(ctx, fmt.Sprintf("token:code:%v", code)); err != nil {
		return err
	}
	if _, err := Postgres().
		ModelContext(ctx, new(Token)).
		Where("code = ?", code).
		Delete(); err != nil {
		return err
	}
	return nil
}

func (s *OtterTokenStore) RemoveByAccess(ctx context.Context, access string) error {
	if err := RedisCache().Delete(ctx, fmt.Sprintf("token:code:%v", access)); err != nil {
		return err
	}
	if _, err := Postgres().
		ModelContext(ctx, new(Token)).
		Where("access = ?", access).
		Delete(); err != nil {
		return err
	}
	return nil
}

func (s *OtterTokenStore) RemoveByRefresh(ctx context.Context, refresh string) error {
	if err := RedisCache().Delete(ctx, fmt.Sprintf("token:code:%v", refresh)); err != nil {
		return err
	}
	if _, err := Postgres().
		ModelContext(ctx, new(Token)).
		Where("refresh = ?", refresh).
		Delete(); err != nil {
		return err
	}
	return nil
}

func (s *OtterTokenStore) GetByCode(ctx context.Context, code string) (oauth2.TokenInfo, error) {
	token := new(Token)
	if err := RedisCache().Once(&cache.Item{
		Ctx:   ctx,
		Key:   fmt.Sprintf("token:code:%v", code),
		Value: token,
		TTL:   15 * time.Minute,
		Do: func(item *cache.Item) (interface{}, error) {
			return s.selectTokenByCode(ctx, code)
		},
	}); err != nil {
		return nil, err
	}
	return token, nil
}

func (s *OtterTokenStore) GetByAccess(ctx context.Context, access string) (oauth2.TokenInfo, error) {
	token := new(Token)
	if err := RedisCache().Once(&cache.Item{
		Ctx:   ctx,
		Key:   fmt.Sprintf("token:access:%v", access),
		Value: token,
		TTL:   15 * time.Minute,
		Do: func(item *cache.Item) (interface{}, error) {
			return s.selectTokenByAccess(ctx, access)
		},
	}); err != nil {
		return nil, err
	}
	return token, nil
}

func (s *OtterTokenStore) GetByRefresh(ctx context.Context, refresh string) (oauth2.TokenInfo, error) {
	token := new(Token)
	if err := RedisCache().Once(&cache.Item{
		Ctx:   ctx,
		Key:   fmt.Sprintf("token:refresh:%v", refresh),
		Value: token,
		TTL:   15 * time.Minute,
		Do: func(item *cache.Item) (interface{}, error) {
			return s.selectTokenByRefresh(ctx, refresh)
		},
	}); err != nil {
		return nil, err
	}
	return token, nil
}

func (s *OtterTokenStore) selectTokenByCode(ctx context.Context, code string) (*Token, error) {
	token := new(Token)
	if err := Postgres().
		ModelContext(ctx, token).
		Where("code = ?", code).
		Select(); err != nil {
		return nil, err
	}
	return token, nil
}

func (s *OtterTokenStore) selectTokenByAccess(ctx context.Context, access string) (*Token, error) {
	token := new(Token)
	if err := Postgres().
		ModelContext(ctx, token).
		Where("access = ?", access).
		Select(); err != nil {
		return nil, err
	}
	return token, nil
}

func (s *OtterTokenStore) selectTokenByRefresh(ctx context.Context, refresh string) (*Token, error) {
	token := new(Token)
	if err := Postgres().
		ModelContext(ctx, token).
		Where("refresh = ?", refresh).
		Select(); err != nil {
		return nil, err
	}
	return token, nil
}
