package observer

import (
	"fmt"
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
