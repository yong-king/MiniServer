package data

import (
	"errors"
	"helloworld/internal/biz"
	"helloworld/internal/conf"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewTodoRepo, NewDB)

// Data .
type Data struct {
	// TODO wrapped database client
	db *gorm.DB
}

// NewData .
func NewData(db *gorm.DB, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{db: db}, cleanup, nil
}

func NewDB(c *conf.Data) (*gorm.DB, error) {
	if c.Database.Dsn == "" {
		return nil, errors.New("database need dsn")
	}
	switch strings.ToLower(c.Database.Driver) {
	case "mysql":
		db, err := gorm.Open(mysql.Open(c.Database.Dsn), &gorm.Config{})
		if err != nil {
			return nil, err
		}
		err = db.AutoMigrate(&biz.Todo{})
		if err != nil {
			return nil, err
		}
		return db, nil
	case "sqlite":
		db, err := gorm.Open(sqlserver.Open(c.Database.Dsn), &gorm.Config{})
		if err != nil {
			return nil, err
		}
		err = db.AutoMigrate(&biz.Todo{})
		if err != nil {
			return nil, err
		}
		return db, nil
	}

	return nil, errors.New("invalid driver")
}
