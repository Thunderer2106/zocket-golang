

package db

import (
	"context"
	"fmt"
	"log"
	"zocket/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/sirupsen/logrus"
)

type LogrusAdapter struct {
	logger *logrus.Logger
}

func (l *LogrusAdapter) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]any) {
	entry := l.logger.WithFields(logrus.Fields(data))

	switch level {
	case tracelog.LogLevelTrace, tracelog.LogLevelDebug:
		entry.Debug(msg)
	case tracelog.LogLevelInfo:
		entry.Info(msg)
	case tracelog.LogLevelWarn:
		entry.Warn(msg)
	case tracelog.LogLevelError:
		entry.Error(msg)
	}
}

var DB *pgxpool.Pool

func ConnectToDB(cfg *config.Config) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("Unable to parse database configuration: %v", err)
	}


	logrusLogger := logrus.New()
	logrusLogger.SetLevel(logrus.DebugLevel) 
	config.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   &LogrusAdapter{logger: logrusLogger},
		LogLevel: tracelog.LogLevelDebug,
	}


	DB, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	log.Println("Connected to PostgreSQL database.")
}
