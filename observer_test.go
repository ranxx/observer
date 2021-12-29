package observer

import (
	"fmt"
	"reflect"
	"testing"
)

// topic, notify, close
type my struct {
	event string `topic:"xxx" notice:"Sayy"`
}

func (m *my) Sayy(name string, age int) {
	fmt.Printf("name:%s age:%d\n", name, age)
}

func TestCheckObserver(t *testing.T) {
	tv, topic, notice := checkObserver(&my{}, "event")
	fmt.Println(tv, topic, notice)
}

func TestSubscribe(t *testing.T) {
	obs := NewObserver()
	obs.Subscribe(&my{})
	obs.Subscribe(&my{})
	obs.Publish("xxx", "Axing", 1000)
	obs.Close()
}

func TestUnsubscribe(t *testing.T) {
	obs := NewObserver()
	obs.Subscribe(&my{})
	obs.Unsubscribe(&my{})
	obs.Publish("xxx", "Axing", 1000)
	obs.Close()
}

type person struct {
	event string `topic:"pain"`
}

func (p *person) Function() string {
	return "Say"
}

func (p *person) Say(name string) {
	fmt.Printf("name:%s 疼痛\n", name)
}

func TestObserver(t *testing.T) {
	obs := NewObserver()
	obs.Subscribe(&my{})
	obs.Subscribe(&person{})
	obs.Publish("xxx", "Axing", 1000)
	obs.Publish("pain", "小明")
	obs.Close()
}

func topicFunc(fc interface{}) {
	fmt.Println(reflect.TypeOf(fc).Kind())
	reflect.ValueOf(fc).Call(nil)
}

func TestFunc(t *testing.T) {
	obs := NewObserver()
	obs.Subscribe(&my{})
	obs.Subscribe(&person{})
	obs.SubscribeByTopicFunc("topic_func", func(v interface{}) {
		fmt.Println("一个", v)
	})
	obs.SubscribeByTopicFunc("topic_func", func(v interface{}) {
		fmt.Println("两个个", v)
	})
	obs.Publish("xxx", "Axing", 1000)
	obs.Publish("pain", "小明")
	obs.Publish("topic_func", &my{})
	obs.Publish("topic_func", &person{})
	obs.Close()
}
