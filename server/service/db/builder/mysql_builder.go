package builder

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"space.online.shop.web.server/service/db/internal"
	"space.online.shop.web.server/util/logger"
)

// = mysqlDbBuilder
type mysqlDbBuilder struct {
	connCfg         *internal.ConnectionCfg
	driverCfg       *mysql.Config
	gormCfg         *gorm.Config
	connMaxIdleTime time.Duration
	maxIdleConns    int
	maxOpenConns    int
}

func (b *mysqlDbBuilder) withOptions(opts ...mysqlDbBuilderOption) {
	for _, opt := range opts {
		opt.apply(b)
	}
}

func (b *mysqlDbBuilder) getDriverConfig() *mysql.Config {
	if b.driverCfg == nil {
		return b.getDefaultDriverConfig()
	}
	return b.driverCfg
}

func (b *mysqlDbBuilder) getDefaultDriverConfig() *mysql.Config {
	return &mysql.Config{
		DSN: b.getDSN(),
	}
}

func (b *mysqlDbBuilder) getDSN() string {
	return b.connCfg.GetDSN()
}

func (b *mysqlDbBuilder) getGormConfig() *gorm.Config {
	if b.gormCfg == nil {
		return b.getDefaultGormConfig()
	}
	return b.gormCfg
}

func (b *mysqlDbBuilder) getDefaultGormConfig() *gorm.Config {
	return &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}
}

func (b *mysqlDbBuilder) getDialector() gorm.Dialector {
	return mysql.New(*b.getDriverConfig())
}

func (b *mysqlDbBuilder) BuildDB() (*gorm.DB, error) {
	db, err := b.build()
	if err != nil {
		return nil, err
	}
	sqlDB, _ := db.DB()
	sqlDB.SetConnMaxIdleTime(b.connMaxIdleTime)
	sqlDB.SetMaxIdleConns(b.maxIdleConns)
	sqlDB.SetMaxOpenConns(b.maxOpenConns)
	return db, nil

}

func (b *mysqlDbBuilder) build() (*gorm.DB, error) {
	db, err := gorm.Open(b.getDialector(), b.getGormConfig())
	if err != nil {
		return nil, err
	}
	_, err = db.DB()
	if err != nil {
		logger.STD.Error("failed to build mysqlDbBuilder, err: %s", err)
		return nil, err
	}
	return db, nil
}

// = mysqlDbBuilderOption

type mysqlDbBuilderOption interface {
	apply(*mysqlDbBuilder)
}

type mysqlBuilderOptionFunc func(*mysqlDbBuilder)

func (fn mysqlBuilderOptionFunc) apply(b *mysqlDbBuilder) {
	fn(b)
}

// = export functions

func NewMysqlDbBuilder(ip, port, database string, opts ...mysqlDbBuilderOption) *mysqlDbBuilder {
	b := &mysqlDbBuilder{
		connCfg:         internal.NewConnectionCfg(ip, port, database),
		maxIdleConns:    10,
		maxOpenConns:    100,
		connMaxIdleTime: 10 * time.Minute,
	}
	b.withOptions(opts...)
	return b
}

func WithCharset(charset string) mysqlDbBuilderOption {
	return mysqlBuilderOptionFunc(func(b *mysqlDbBuilder) {
		b.connCfg.ExtraConncetionCfg.Charset = charset
	})
}

func WithParseTime(parseTime string) mysqlDbBuilderOption {
	return mysqlBuilderOptionFunc(func(b *mysqlDbBuilder) {
		b.connCfg.ExtraConncetionCfg.ParseTime = parseTime
	})
}

func WithLoc(loc string) mysqlDbBuilderOption {
	return mysqlBuilderOptionFunc(func(b *mysqlDbBuilder) {
		b.connCfg.ExtraConncetionCfg.Loc = loc
	})
}

func WithTimeout(t string) mysqlDbBuilderOption {
	return mysqlBuilderOptionFunc(func(b *mysqlDbBuilder) {
		b.connCfg.ExtraConncetionCfg.Timeout = t
	})
}

func WithWriteTimeout(t string) mysqlDbBuilderOption {
	return mysqlBuilderOptionFunc(func(b *mysqlDbBuilder) {
		b.connCfg.ExtraConncetionCfg.WriteTimeout = t
	})
}

func WithReadTimeout(t string) mysqlDbBuilderOption {
	return mysqlBuilderOptionFunc(func(b *mysqlDbBuilder) {
		b.connCfg.ExtraConncetionCfg.ReadTimeout = t
	})
}

func WithDriverCfg(cfg *mysql.Config) mysqlDbBuilderOption {
	return mysqlBuilderOptionFunc(func(b *mysqlDbBuilder) {
		b.driverCfg = cfg
	})
}

func WithGormCfg(cfg *gorm.Config) mysqlDbBuilderOption {
	return mysqlBuilderOptionFunc(func(b *mysqlDbBuilder) {
		b.gormCfg = cfg
	})
}

func WithMaxIdleConns(n int) mysqlDbBuilderOption {
	return mysqlBuilderOptionFunc(func(b *mysqlDbBuilder) {
		b.maxIdleConns = n
	})
}

func WithMaxOpenConns(n int) mysqlDbBuilderOption {
	return mysqlBuilderOptionFunc(func(b *mysqlDbBuilder) {
		b.maxOpenConns = n
	})
}

func WithConnMaxIdleTime(t time.Duration) mysqlDbBuilderOption {
	return mysqlBuilderOptionFunc(func(b *mysqlDbBuilder) {
		b.connMaxIdleTime = t
	})
}
