# observer

设计模式之观察者模式

## 快速上手

可以查看[observer_test](./observer_test.go)，了解怎么使用

```go
package main

import (
	"fmt"
	"github.com/ranxx/observer"
)

func main() {
	obs := observer.NewObserver[int64]()
	obs.SubscribeByTopicFunc("func_arg_slice", func(a string, b ...int) {
		fmt.Println(a, b)
	})
	obs.Publish("func_arg_slice", -99, "axing", 1, 2, 3, 4, 5)
	obs.Wait()
}
```

### 函数方式

```go
package main

import (
	"fmt"
	"github.com/ranxx/observer"
)

func main() {
	obs := observer.NewObserver[interface{}]()
	obs.SubscribeByTopicFunc("func_arg_slice", func(a string, b ...int) {
		fmt.Println(a, b)
	})
	obs.SyncPublish("func_arg_slice", nil, "axing", 1, 2, 3, 4, 5)
	obs.SyncPublishWithRet("func_arg_slice", nil, func(){
		fmt.Println("阿星 执行完毕")
	}, "阿星", 1, 2, 3, 4, 5)
}
```

### 结构体方式

```go
package main

import (
	"fmt"
	"github.com/ranxx/observer"
)

type person struct {
	event string `topic:"pain" notice:"Say"`
	name  string `json:"name"`
}

func (p *person) Say() {
	fmt.Printf("name:%s 疼痛\n", p.name)
}

func main() {
	obs := observer.NewObserver[string]()
	obs.Subscribe(&person{name:"小明"})
	obs.Subscribe(&person{name:"小红"})
	obs.Publish("pain", "common")
	obs.Wait()
}
```