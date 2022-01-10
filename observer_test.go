package observer

import (
	"fmt"
	"reflect"
	"testing"
)

// topic, notify, close
type my struct {
	event string `topic:"xxx" notice:"Sayy"`
	name  string `json:"name"`
}

func (m *my) Sayy(name string, age int) {
	if len(name) == 0 {
		name = m.name
	}
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
	obs.Wait()
}

func TestUnsubscribe(t *testing.T) {
	obs := NewObserver()
	cancel := obs.Subscribe(&my{name: "小明"})
	_ = obs.Subscribe(&my{name: "小红"})
	_ = obs.Subscribe(&my{name: "小刚"})
	cancel()
	obs.Publish("xxx", "", 1000)
	obs.Wait()
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
	obs.Wait()
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
	obs.SubscribeByTopicFunc("topic_func1", func(v interface{}, name string, age int) {
		fmt.Println("两个个", v, name, age)
	})
	obs.Publish("xxx", "Axing", 1000)
	obs.Publish("pain", "小明")
	obs.Publish("topic_func", &my{})
	obs.Publish("topic_func", &person{})
	obs.Publish("topic_func1", &person{}, "ainxg", 24)
	obs.Wait()
}

func TestPublishWithRet(t *testing.T) {
	obs := NewObserver()
	obs.SubscribeByTopicFunc("topic_func_1_ret", func(v interface{}) (string, error) {
		fmt.Println("一个", v)
		return "第一个", fmt.Errorf("一个错误")
	})
	obs.SubscribeByTopicFunc("topic_func_1_ret", func(v interface{}) (string, error) {
		fmt.Println("二个", v)
		return "第二个", nil
	})

	obs.PublishWithRet("topic_func_1_ret", func(v string, e error) {
		if e != nil {
			fmt.Println("返回报错", e)
			return
		}
		fmt.Println("返回", v)
	}, "鸡蛋")
	obs.Wait()
}

func TestSyncPublish(t *testing.T) {
	obs := NewObserver()
	obs.SubscribeByTopicFunc("sync_topic_func_ret_1", func(v interface{}) (string, error) {
		fmt.Println("一个", v)
		return "第一个", fmt.Errorf("一个错误")
	})
	obs.SubscribeByTopicFunc("sync_topic_func_ret_1", func(v interface{}) (string, error) {
		fmt.Println("二个", v)
		return "第二个", nil
	})

	obs.SyncPublish("sync_topic_func_ret_1", "鸡蛋")
}

func TestSyncPublishWithRet(t *testing.T) {
	obs := NewObserver()
	obs.SubscribeByTopicFunc("sync_topic_func_ret_1", func(v interface{}) (string, error) {
		fmt.Println("一个", v)
		return "第一个", fmt.Errorf("一个错误")
	})

	obs.SubscribeByTopicFunc("sync_topic_func_ret_1", func(v interface{}) (string, error) {
		fmt.Println("二个", v)
		return "第二个", nil
	})

	obs.SubscribeByTopicFunc("sync_topis_func_ret_2", func(v interface{}) {
		fmt.Println("第三个")
	})

	obs.SyncPublishWithRet("sync_topic_func_ret_1", func(v string, err error) {
		if err != nil {
			fmt.Println("返回报错", err)
			return
		}
		fmt.Println("返回", v)

	}, "鸡蛋")

	obs.SyncPublishWithRet("sync_topis_func_ret_2", func() {
		fmt.Println("没有返回值")
	}, "鸡蛋")
}
