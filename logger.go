package zapion

import (
	"sync"

	"github.com/pion/logging"
	"go.uber.org/zap"
)

type logger struct {
	*zap.SugaredLogger
	wrapped     *zap.SugaredLogger
	enableTrace bool
}

func (l *logger) Trace(msg string) {
	if l.enableTrace {
		l.wrapped.Debug(msg)
	}
}

func (l *logger) Tracef(format string, args ...interface{}) {
	if l.enableTrace {
		l.wrapped.Debugf(format, args...)
	}
}

func (l *logger) Debug(msg string) {
	l.wrapped.Debug(msg)
}

func (l *logger) Info(msg string) {
	l.wrapped.Info(msg)
}

func (l *logger) Warn(msg string) {
	l.wrapped.Warn(msg)
}

func (l *logger) Error(msg string) {
	l.wrapped.Error(msg)
}

// ZapFactory is a logger factory backended by zap logger.
type ZapFactory struct {
	BaseLogger  *zap.Logger
	EnableTrace bool

	mu      sync.Mutex
	loggers []*logger
}

// NewLogger creates new scoped logger.
func (f *ZapFactory) NewLogger(scope string) logging.LeveledLogger {
	f.mu.Lock()
	defer f.mu.Unlock()

	named := f.BaseLogger.Named(scope)
	l := &logger{
		SugaredLogger: named.Sugar(),
		wrapped:       named.WithOptions(zap.AddCallerSkip(1)).Sugar(),
		enableTrace:   f.EnableTrace,
	}
	f.loggers = append(f.loggers, l)
	return l
}

// SyncAll calls Sync() method of all child loggers.
// It is recommended to call this before exiting the program to
// ensure never loosing buffered log.
func (f *ZapFactory) SyncAll() {
	f.mu.Lock()
	defer f.mu.Unlock()

	for _, l := range f.loggers {
		_ = l.SugaredLogger.Sync()
	}
}
