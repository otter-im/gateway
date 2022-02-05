package app

import (
	"github.com/go-pg/pg/v10"
	"github.com/otter-im/auth/internal/config"
	"sync"
)

var (
	pgMainOnce sync.Once
	pgMain     *pg.DB
)

func Postgres() *pg.DB {
	pgMainOnce.Do(func() {
		options := &pg.Options{
			Addr:     config.PostgresAddress(),
			User:     config.PostgresUser(),
			Password: config.PostgresPassword(),
			Database: config.PostgresDatabase(),
		}

		pgMain = pg.Connect(options)
		AddExitHook(func() error {
			if err := pgMain.Close(); err != nil {
				return err
			}
			return nil
		})
	})
	return pgMain
}