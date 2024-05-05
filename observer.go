package observer

import (
	"reflect"
	"sync"
)

const (
	event = "event"
)

// Observer 观察者 模式
type Observer[E any] interface {
	Subscriber
	Publisher[E]
	Wait()
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
	// 订阅, 去重
	//
	// 返回的func是为 取消订阅
	Subscribe(observer interface{}) func()

	// // 取消订阅
	// Unsubscribe(observer interface{})

	// 订阅 topic 函数, 不去重
	//
	// 返回的func是为 取消订阅
	SubscribeByTopicFunc(topic string, fc interface{}) func()
}

// Publisher 发布者
type Publisher[E any] interface {
	// 异步推送
	Publish(topic string, extra E, args ...interface{})

	// 异步推送,监听返回值
	//
	// @params: retfc 其参数类型和顺序与消费者保持一致, 如没有返回值可为空
	PublishWithRet(topic string, extra E, retfc interface{}, args ...interface{})

	// 同步推送
	SyncPublish(topic string, extra E, args ...interface{})

	// 同步推送,监听返回值
	//
	// @params: retfc 其参数类型和顺序与消费者保持一致, 如没有返回值可为空
	SyncPublishWithRet(topic string, extra E, retfc interface{}, args ...interface{})
}

// Topic 主题名
type Topic interface {
	Topic() string
}

// Function 通知函数
type Function interface {
	Function() string
}

type observer[E any] struct {
	handler *syncHandler
	wait    sync.WaitGroup
	opt     *Options[E]
}

// NewObserver Observer
func NewObserver[E any](opts ...Option[E]) Observer[E] {
	opt := _defaultOpt[E]()
	for _, v := range opts {
		v(opt)
	}
	return &observer[E]{
		handler: &syncHandler{
			rwlock: new(sync.RWMutex),
			m:      map[string]handlers{},
		},
		opt: opt,
	}
}

func (o *observer[E]) Subscribe(observer interface{}) func() {
	t, topic, fc := checkObserver(observer, o.opt.eventField)

	h := newHandler(t, observer, fc)

	o.handler.Append(topic, true, h)
	return func() {
		o.handler.Del(topic, h)
	}
}

// func (o *observer[E]) Unsubscribe(observer interface{}) {
// 	t, topic, fc := checkObserver(observer, event)

// 	o.handler.Del(topic, newHandler(t, observer, fc))
// }

func (o *observer[E]) SubscribeByTopicFunc(topic string, fc interface{}) func() {
	t, v := checkFunc(fc)

	h := newHandler(t, fc, v)

	o.handler.Append(topic, false, h)
	return func() {
		o.handler.Del(topic, h)
	}
}

func (o *observer[E]) Publish(topic string, extra E, args ...interface{}) {
	o.publish(false, topic, extra, nil, args...)
}

func (o *observer[E]) PublishWithRet(topic string, extra E, retfc interface{}, args ...interface{}) {
	o.publish(false, topic, extra, retfc, args...)
}

func (o *observer[E]) SyncPublish(topic string, extra E, args ...interface{}) {
	o.publish(true, topic, extra, nil, args...)
}

func (o *observer[E]) SyncPublishWithRet(topic string, extra E, retfc interface{}, args ...interface{}) {
	o.publish(true, topic, extra, retfc, args...)
}

func (o *observer[E]) publish(sync bool, topic string, extra E, retfc interface{}, args ...interface{}) {
	fc := reflect.Zero(reflect.TypeOf(0))
	if retfc != nil {
		_, fc = checkFunc(retfc)
	}

	handlers := o.handler.Get(topic)
	params := []reflect.Value{}

	for _, arg := range args {
		params = append(params, reflect.ValueOf(arg))
	}

	for _, handler := range handlers {
		o.onNotice(sync, topic, extra, handler, fc, params, args)
	}
}

func (o *observer[E]) onNotice(sync bool, topic string, extra E, handler *handler, fc reflect.Value, params []reflect.Value, args ...interface{}) {
	h := func() {
		defer o.opt.recover(topic, extra, handler.callback, args)
		defer o.wait.Done()

		ret := handler.Call(params...)
		if !fc.IsZero() {
			fc.Call(ret)
		}
	}
	o.wait.Add(1)
	if sync {
		o.opt.syncExec(topic, extra, h)
		return
	}
	o.opt.asyncExec(topic, extra, h)
}

func (o *observer[E]) Wait() {
	o.wait.Wait()
}
