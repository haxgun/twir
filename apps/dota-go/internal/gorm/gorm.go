package gorm

import (
	"context"
	"time"

	cfg "github.com/satont/twir/libs/config"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func New(config cfg.Config, lc fx.Lifecycle, logger *logrus.Logger) (*gorm.DB, error) {
	db, err := gorm.Open(
		postgres.Open(config.DatabaseUrl), &gorm.Config{
			Logger: gormLogger.Default.LogMode(gormLogger.Silent),
		},
	)
	if err != nil {
		return nil, err
	}
	d, _ := db.DB()
	d.SetMaxOpenConns(5)
	d.SetConnMaxIdleTime(1 * time.Minute)

	lc.Append(
		fx.Hook{
			OnStop: func(ctx context.Context) error {
				return d.Close()
			},
		},
	)

	logger.Infoln("Postgres connected")

	return db, nil
}
