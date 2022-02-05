package logger

import (
	"context"
	"github.com/tryfix/log"
)

var Log log.Logger

func Init(ctx context.Context) {
	Log = log.Constructor.Log(
		log.WithColors(config.ColorsEnabled),
		log.WithLevel(log.Level(config.LogLevel)),
		log.WithFilePath(config.FilePath),
	)

	Log.InfoContext(ctx, `logger initialized`)
}