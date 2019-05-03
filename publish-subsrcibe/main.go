package main

import (
	"errors"
	"fmt"
	"time"
)

func main() {

}

type Subscription struct {
	topics []chan<- string
}

func (s Subscription) Publish(msg string) error {
	if len(s.topics) == 0 {
		return errors.New("no topics")
	}
	for _, ch := range s.topics {
		select {
		case ch <- msg:
		case <-time.After(1 * time.Second):
			fmt.Println("this topic is timeout")
		}
	}
	return nil
}

func (s *Subscription) AddTopic() (<-chan string, error) {
	ch := make(chan string)
	s.topics = append(s.topics, ch)
	return ch, nil
}

type Topic struct {
	listen_ch <-chan string
}

func (t *Topic) Subscribe(s Subscription) error {
	t.listen_ch, _ = s.AddTopic()
	return nil
}
