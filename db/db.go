package db

import (
	"github.com/infinytum/injector"
	"gorm.io/gorm"
)

type cachedDb struct {
	db *gorm.DB
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
