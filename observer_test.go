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
	obs := NewObserver[interface{}]()
	obs.Subscribe(&my{})
	obs.Subscribe(&my{})
	obs.Publish("xxx", nil, "Axing", 1000)
	obs.Wait()
}

func TestUnsubscribe(t *testing.T) {
	obs := NewObserver[int64]()
	cancel := obs.Subscribe(&my{name: "小明"})
	_ = obs.Subscribe(&my{name: "小红"})
	_ = obs.Subscribe(&my{name: "小刚"})
	cancel()
	obs.Publish("xxx", 0, "", 1000)
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
	obs := NewObserver[int64]()
	obs.Subscribe(&my{})
	obs.Subscribe(&person{})
	obs.Publish("xxx", 0, "Axing", 1000)
	obs.Publish("pain", 2, "小明")
	obs.Wait()
}

func topicFunc(fc interface{}) {
	fmt.Println(reflect.TypeOf(fc).Kind())
	reflect.ValueOf(fc).Call(nil)
}

func TestFunc(t *testing.T) {
	obs := NewObserver[int64]()
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
	obs.Publish("xxx", 0, "Axing", 1000)
	obs.Publish("pain", 1, "小明")
	obs.Publish("topic_func", 2, &my{})
	obs.Publish("topic_func", 3, &person{})
	obs.Publish("topic_func1", 4, &person{}, "ainxg", 24)
	obs.Wait()
}

func TestPublishWithRet(t *testing.T) {
	obs := NewObserver[string]()
	obs.SubscribeByTopicFunc("topic_func_1_ret", func(v interface{}) (string, error) {
		fmt.Println("一个", v)
		return "第一个", fmt.Errorf("一个错误")
	})
	obs.SubscribeByTopicFunc("topic_func_1_ret", func(v interface{}) (string, error) {
		fmt.Println("二个", v)
		return "第二个", nil
	})

	obs.PublishWithRet("topic_func_1_ret", "1", func(v string, e error) {
		if e != nil {
			fmt.Println("返回报错", e)
			return
		}
		fmt.Println("返回", v)
	}, "鸡蛋")
	obs.Wait()
}

func TestSyncPublish(t *testing.T) {
	obs := NewObserver[int64]()
	obs.SubscribeByTopicFunc("sync_topic_func_ret_1", func(v interface{}) (string, error) {
		fmt.Println("一个", v)
		return "第一个", fmt.Errorf("一个错误")
	})
	obs.SubscribeByTopicFunc("sync_topic_func_ret_1", func(v interface{}) (string, error) {
		fmt.Println("二个", v)
		return "第二个", nil
	})

	obs.SyncPublish("sync_topic_func_ret_1", 1, "鸡蛋")
}

func TestSyncPublishWithRet(t *testing.T) {
	obs := NewObserver[int64]()
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

	obs.SyncPublishWithRet("sync_topic_func_ret_1", 1, func(v string, err error) {
		if err != nil {
			fmt.Println("返回报错", err)
			return
		}
		fmt.Println("返回", v)

	}, "鸡蛋")

	obs.SyncPublishWithRet("sync_topis_func_ret_2", 1, func() {
		fmt.Println("没有返回值")
	}, "鸡蛋")
}

func TestFuncSlice(t *testing.T) {
	obs := NewObserver[string]()
	obs.SubscribeByTopicFunc("func_arg_slice", func(a string, b ...int) {
		fmt.Println(a, b)
	})
	obs.Publish("func_arg_slice", "axing-extra", "axing", 1, 2, 3, 4, 5)
	obs.Wait()
}

func TestStruct(t *testing.T) {
	obs := NewObserver[string]()
	obs.Subscribe(struct {
		event string `topic:"topic"`
	}{})
	obs.Wait()
}

type structV2 struct {
	event string `notice:"Handler"`
	Name  string
}

func (s *structV2) Topic() string {
	return "structV2"
}

func (s *structV2) Handler() {
	fmt.Println(s.Name)
}

func TestStructV2(t *testing.T) {
	obs := NewObserver[string]()

	// Topic
	obs.Subscribe(&structV2{Name: "小明"})
	obs.SubscribeByTopicFunc("structV2", func() {
		fmt.Println("topic-func 执行了")
	})
	obs.SubscribeByTopicFunc("structV2", func() {
		fmt.Println("topic-func 执行了 - 1")
	})

	obs.Publish("structV2", "ssss")
	obs.SyncPublish("structV2", "aaaa")
	obs.Wait()
}

func TestStructV3(t *testing.T) {
	obs := NewObserver(WithSyncExec(func(topic string, extra string, fn func()) {
		fmt.Println(topic, extra, "开始执行")
		fn()
		fmt.Println(topic, extra, "执行完毕")
	}), WithAsyncExec(func(topic string, extra string, fn func()) {
		fmt.Println(topic, extra, "开始执行")
		fn()
		fmt.Println(topic, extra, "执行完毕")
	}))

	// Topic
	obs.Subscribe(&structV2{Name: "小明"})
	obs.SubscribeByTopicFunc("structV2", func() {
		fmt.Println("topic-func 执行了")
	})
	obs.SubscribeByTopicFunc("structV2", func() {
		fmt.Println("topic-func 执行了 - 1")
	})

	obs.Publish("structV2", "ssss")
	obs.SyncPublish("structV2", "aaaa")
	obs.Wait()
}
