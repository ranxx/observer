# observer

设计模式之观察者模式

## 快速上手

```go
import "github.com/ranxx/observer"

type person struct {
	event string `topic:"pain"`
}

func (p *person) Function() string {
	return "Say"
}

func (p *person) Say(name string) {
	fmt.Printf("name:%s 疼痛\n", name)
}

func main() {
	obs := observer.NewObserver()
	obs.Subscribe(&person{})
	obs.Publish("pain", "小明")
	obs.Close()
}v
```

