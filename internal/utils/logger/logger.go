package logger

import (
	"sync"

	"go.uber.org/zap"
)

var (
	once sync.Once
	l    *zap.Logger
)

// Init initializes the global logger. It only takes effect once;
// subsequent calls are no-ops. It is also called lazily on first use.
func Init() {
	once.Do(func() {
		var err error
		l, err = zap.NewDevelopment()
		if err != nil {
			panic(err)
		}
	})
}

// L returns the global *zap.Logger, initializing it on first use.
func L() *zap.Logger {
	Init()
	return l
}

func Debug(msg string, fields ...zap.Field) { L().Debug(msg, fields...) }
func Info(msg string, fields ...zap.Field)  { L().Info(msg, fields...) }
func Warn(msg string, fields ...zap.Field)  { L().Warn(msg, fields...) }
func Error(msg string, fields ...zap.Field) { L().Error(msg, fields...) }
