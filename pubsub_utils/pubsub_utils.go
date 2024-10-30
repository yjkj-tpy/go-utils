package pubsub_utils

type EsPubSub struct {
	subscribers map[string][]chan interface{}
}

// NewEsPubSub 获取一个在本地内存的消费订阅，注意不是单例的
func NewEsPubSub() *EsPubSub {
	return &EsPubSub{
		subscribers: make(map[string][]chan interface{}),
	}
}

func (ps *EsPubSub) Subscribe(topic string, size int) chan interface{} {
	ch := make(chan interface{}, size)
	ps.subscribers[topic] = append(ps.subscribers[topic], ch)
	return ch
}

func (ps *EsPubSub) Publish(topic string, message interface{}) {
	if subs, ok := ps.subscribers[topic]; ok {
		for _, sub := range subs {
			go func(ch chan interface{}) {
				ch <- message
			}(sub)
		}
	}
}
