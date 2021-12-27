package observer

import (
	"fmt"
	"reflect"
	"sync"
)

const (
	event = "event"
)

// Observer 观察者 模式
type Observer interface {
	Subscriber
	Publisher
	Close()
}

// Subscriber 订阅者
//
// param observer must be a reflect.Struct
//
// observer can provide Topic or Function interface
//
// or probide event field.
//
// for example:
/*
type person struct {
	event string `topic:"pain" notice:"Say"`
}

func(p*person) Say() {
	// do sth
}
*/
type Subscriber interface {
	// 订阅
	Subscribe(observer interface{})
	// 取消订阅
	Unsubscribe(observer interface{})
}

// Publisher 发布者
type Publisher interface {
	Publish(topic string, args ...interface{})
}

// Topic 主题名
type Topic interface {
	Topic() string
}

// Function 通知函数
type Function interface {
	Function() string
}

func checkObserver(observer interface{}, fieldName string) (reflect.Type, string, string) {
	t, topic, function := checkObserverForInterface(observer)

	if len(topic) > 0 && len(function) > 0 {
		return t, topic, function
	}

	// 字段名
	field, ok := t.FieldByName(fieldName)
	if !ok {
		panic(fmt.Sprintf("%s is no %s field", t.String(), fieldName))
	}

	if len(topic) <= 0 {
		topic, ok = field.Tag.Lookup("topic")
	}

	if !ok {
		panic(fmt.Sprintf("%s.%s field is no topic in the tag", t.String(), fieldName))
	}

	if len(function) <= 0 {
		function, ok = field.Tag.Lookup("notice")
	}

	if !ok {
		panic(fmt.Sprintf("%s.%s field is no notice in the tag", t.String(), fieldName))
	}

	return t, topic, function
}

func checkObserverForInterface(observer interface{}) (reflect.Type, string, string) {
	t := reflect.TypeOf(observer)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		panic(fmt.Sprintf("%s is not reflect.Struct", t.String()))
	}

	var topic, function string

	topicf, ok := observer.(Topic)
	if ok {
		topic = topicf.Topic()
	}

	functionf, ok := observer.(Function)
	if ok {
		function = functionf.Function()
	}

	return t, topic, function
}

type observer struct {
	handler *syncHandler
	wait    sync.WaitGroup
}

// NewObserver Observer
func NewObserver() Observer {
	return &observer{
		handler: &syncHandler{
			rwlock: new(sync.RWMutex),
			m:      map[string]handlers{},
		},
	}
}

func (o *observer) Subscribe(observer interface{}) {
	t, topic, function := checkObserver(observer, event)

	o.handler.Append(topic, newHandler(t, observer, function))
}

func (o *observer) Unsubscribe(observer interface{}) {
	t, topic, function := checkObserver(observer, event)

	o.handler.Del(topic, newHandler(t, observer, function))
}

func (o *observer) Publish(topic string, args ...interface{}) {
	handlers := o.handler.Get(topic)
	params := []reflect.Value{}

	for _, arg := range args {
		params = append(params, reflect.ValueOf(arg))
	}

	for _, handler := range handlers {
		o.onNotice(handler, params)
	}
}

func (o *observer) onNotice(handler *handler, params []reflect.Value) {
	o.wait.Add(1)
	go func() {
		defer o.wait.Done()
		handler.Call(params...)
	}()
}

func (o *observer) Close() {
	o.wait.Wait()
}
