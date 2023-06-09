//nolint:gocritic // here it is fine
package logger

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var _ Log = &logger{}

type Log interface {
	RawJSON(key string, value []byte) Log
	Duration(key string, duration time.Duration) Log
	String(key, value string) Log
	Err(error) Log
	Int(key string, value int) Log
	Bool(key string, value bool) Log
	Msg(msg string)
}

type logger struct {
	zeroLog *zerolog.Event
}

// Err implements Log
func (l *logger) Err(err error) Log {
	l.zeroLog.Err(err)
	return l
}

func (l *logger) RawJSON(key string, value []byte) Log {
	l.zeroLog.RawJSON(key, value)
	return l
}

func (l *logger) Bool(key string, value bool) Log {
	l.zeroLog.Bool(key, value)
	return l
}

func (l *logger) Duration(key string, duration time.Duration) Log {
	l.zeroLog.Dur(key, duration)
	return l
}

func (l *logger) Int(key string, value int) Log {
	l.zeroLog.Int(key, value)
	return l
}

func (l *logger) Msg(msg string) {
	l.zeroLog.Msg(msg)
}

func (l *logger) String(key, value string) Log {
	l.zeroLog.Str(key, value)

	return l
}

func Fatal(_ context.Context, err error) Log {
	log := log.Fatal().Err(err)
	return &logger{zeroLog: log}
}

func Error(_ context.Context, err error) Log {
	log := log.Err(err)
	return &logger{zeroLog: log}
}

func Info(_ context.Context) Log {
	log := log.Info()
	return &logger{zeroLog: log}
}

func Debug(_ context.Context) Log {
	log := log.Debug()
	return &logger{zeroLog: log}
}
