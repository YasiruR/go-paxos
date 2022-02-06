package logger

import (
	"context"
	"errors"
	"fmt"
	"github.com/tryfix/log"
	"runtime"
)

func Init(ctx context.Context) log.Logger {
	logger := log.Constructor.Log(
		log.WithColors(config.ColorsEnabled),
		log.WithLevel(log.Level(config.LogLevel)),
		log.WithFilePath(config.FilePath),
	)

	logger.InfoContext(ctx, `logger initialized`)
	return logger
}

func ErrorWithLine(err error) error {
	_, file, line, _ := runtime.Caller(1)
	return errors.New(fmt.Sprintf(`%s - %s:%d`, err.Error(), file, line))
}
