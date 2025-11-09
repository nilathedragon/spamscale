package db

import (
	"github.com/go-mojito/mojito"
	"github.com/go-mojito/mojito/pkg/cache"
	"github.com/infinytum/injector"
	"gorm.io/gorm"
)

type cachedDb struct {
	db    *gorm.DB
	cache cache.Cache
}

func (c *cachedDb) getDB() *gorm.DB {
	if c.db == nil {
		db, err := injector.Inject[*gorm.DB]()
		if err != nil {
			panic(err)
		}
		c.db = db
	}
	return c.db
}

func (c *cachedDb) getCache() cache.Cache {
	if c.cache == nil {
		c.cache = mojito.DefaultCache()
	}
	return c.cache
}
