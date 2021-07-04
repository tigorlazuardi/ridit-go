package pkg

import (
	"context"

	"github.com/sirupsen/logrus"
)

type M = logrus.Fields

type entryKey struct{}

var ContextEntryKey = entryKey{}

func ContextWithEntry(ctx context.Context, entry *logrus.Entry) context.Context {
	return context.WithValue(ctx, ContextEntryKey, entry)
}

func EntryFromContext(ctx context.Context) *logrus.Entry {
	entry, ok := ctx.Value(ContextEntryKey).(*logrus.Entry)
	if !ok {
		return logrus.NewEntry(logrus.StandardLogger())
	}
	return entry
}

func ContextEntryWithFields(ctx context.Context, fields logrus.Fields) context.Context {
	entry := EntryFromContext(ctx)
	return ContextWithEntry(ctx, entry.WithFields(fields))
}

func ContextEntryWithError(ctx context.Context, err error) context.Context {
	entry := EntryFromContext(ctx)
	return ContextWithEntry(ctx, entry.WithError(err))
}

func ContextLog(ctx context.Context, level logrus.Level, message ...interface{}) {
	entry := EntryFromContext(ctx)
	entry.WithContext(ctx).Log(level, message...)
}

func ContextLogWithFields(ctx context.Context, level logrus.Level, fields logrus.Fields, message ...interface{}) {
	entry := EntryFromContext(ctx)
	entry.WithContext(ctx).WithFields(fields).Log(level, message...)
}

func ContextLogWithError(ctx context.Context, level logrus.Level, err error, message ...interface{}) {
	entry := EntryFromContext(ctx)
	entry.WithContext(ctx).WithError(err).Log(level, message...)
}

func ContextLogWithErrorAndFields(ctx context.Context, level logrus.Level, err error, fields logrus.Fields, message ...interface{}) {
	entry := EntryFromContext(ctx)
	entry.WithContext(ctx).WithFields(fields).WithError(err).Log(level, message...)
}
