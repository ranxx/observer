package observer

import "reflect"

// Options options
type Options[E any] struct {
	eventField string // default event

	// recover
	recover func(topic string, extra E, cb reflect.Value, args ...interface{})

	// exec
	asyncExec func(topic string, extra E, fn func())
	syncExec  func(topic string, extra E, fn func())
}

// Option option
type Option[E any] func(o *Options[E])

func _defaultOpt[E any]() *Options[E] {
	return &Options[E]{
		eventField: event,

		recover: func(topic string, extra E, cb reflect.Value, args ...interface{}) {
			if err := recover(); err != nil {
				panic([]interface{}{topic, extra, cb.Type().PkgPath(), cb.Type().Name(), cb.String(), args, err})
			}
		},

		asyncExec: func(topic string, extra E, fn func()) { go fn() },
		syncExec:  func(topic string, extra E, fn func()) { fn() },
	}
}

// WithEventField event field default event
func WithEventField[E any](event string) Option[E] {
	return func(o *Options[E]) {
		o.eventField = event
	}
}

// WithRecover recover default
func WithRecover[E any](re func(topic string, extra E, cb reflect.Value, args ...interface{})) Option[E] {
	return func(o *Options[E]) {
		o.recover = re
	}
}

// WithAsyncExec async exec
func WithAsyncExec[E any](exec func(topic string, extra E, fn func())) Option[E] {
	return func(o *Options[E]) {
		o.asyncExec = exec
	}
}

// WithSyncExec sync exec
func WithSyncExec[E any](exec func(topic string, extra E, fn func())) Option[E] {
	return func(o *Options[E]) {
		o.syncExec = exec
	}
}
