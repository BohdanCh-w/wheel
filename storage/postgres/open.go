package postgres

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/bohdanch-w/wheel/logger"
)

type DBOption func(*gorm.Config)

func Open(dsn string, opts ...DBOption) (*gorm.DB, error) {
	cfg := gorm.Config{}

	for _, opt := range opts {
		opt(&cfg)
	}

	return gorm.Open(postgres.Open(dsn), &cfg) // nolint: wrapcheck
}

func WithLogLevel(level logger.LogLevel) DBOption {
	return func(config *gorm.Config) {
		dbLogger := gormlogger.Discard

		switch level {
		case logger.Debug:
			dbLogger = gormlogger.Default.LogMode(gormlogger.Info)
		case logger.Info:
			dbLogger = gormlogger.Default.LogMode(gormlogger.Warn)
		case logger.Warn:
			dbLogger = gormlogger.Default.LogMode(gormlogger.Error)
		case logger.Error:
			dbLogger = gormlogger.Default.LogMode(gormlogger.Silent)
		}

		config.Logger = dbLogger
	}
}
