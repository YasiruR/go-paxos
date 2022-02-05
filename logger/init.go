package logger

import (
	"context"
	"errors"
	"fmt"
	"github.com/tryfix/log"
	"runtime"
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

func ErrorWithLine(err error) error {
	_, file, line, _ := runtime.Caller(1)
	return errors.New(fmt.Sprintf(`%s - %s:%d`, err.Error(), file, line))
}
