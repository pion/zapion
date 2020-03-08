package zapion

import (
	"sync"

	"github.com/pion/logging"
	"go.uber.org/zap"
)

type logger struct {
	l           *zap.SugaredLogger
	enableTrace bool
}

func (l *logger) Trace(msg string) {
	if l.enableTrace {
		l.l.Debug(msg)
	}
}

func (l *logger) Tracef(format string, args ...interface{}) {
	if l.enableTrace {
		l.l.Debugf(format, args...)
	}
}

func (l *logger) Debug(msg string) {
	l.l.Debug(msg)
}

func (l *logger) Debugf(format string, args ...interface{}) {
	l.l.Debugf(format, args...)
}

func (l *logger) Info(msg string) {
	l.l.Info(msg)
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.l.Infof(format, args...)
}

func (l *logger) Warn(msg string) {
	l.l.Warn(msg)
}

func (l *logger) Warnf(format string, args ...interface{}) {
	l.l.Warnf(format, args...)
}

func (l *logger) Error(msg string) {
	l.l.Error(msg)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.l.Errorf(format, args...)
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

	l := &logger{
		l: f.BaseLogger.
			WithOptions(zap.AddCallerSkip(1)).
			Named(scope).
			Sugar(),
		enableTrace: f.EnableTrace,
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
		_ = l.l.Sync()
	}
}
