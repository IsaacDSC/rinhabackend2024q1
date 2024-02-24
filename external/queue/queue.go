package queue

import "time"

type Input = []byte

type Consumer interface {
	Consumer(t Input) error
	ConsumerErr(t Retry) error
}

type Retry struct {
	RetrieveError  error
	Quantity       int
	Msg            Input
	TimeoutSeconds int
}

type Event struct {
	maxThread   int
	bufferQueue []Input
	retryQueue  []Retry
	pubTo       Consumer
}

func NewEvent(consumer Consumer) *Event {
	return &Event{
		maxThread: 50,
		pubTo:     consumer,
	}
}

func (e *Event) rmQueue() {
	e.bufferQueue = e.bufferQueue[1:len(e.bufferQueue)]
}

func (e *Event) rmQueueErr() {
	e.retryQueue = e.retryQueue[1:len(e.retryQueue)]
}

func (e *Event) Consume() {
	go e.consumeErr()
	for {
		if len(e.bufferQueue) > 0 {
			err := e.pubTo.Consumer(e.bufferQueue[0])
			if err != nil {
				e.retryQueue = append(e.retryQueue, Retry{
					Quantity:      1,
					RetrieveError: err,
					Msg:           e.bufferQueue[0],
				})
			}
			e.rmQueue()
		}
	}
}

func (e *Event) consumeErr() {
	for {
		if len(e.retryQueue) > 0 {
			if e.retryQueue[0].Quantity <= 5 {
				e.retryQueue[0].TimeoutSeconds += 1
				time.Sleep(time.Second * time.Duration(e.retryQueue[0].TimeoutSeconds))
				if err := e.pubTo.ConsumerErr(e.retryQueue[0]); err != nil {
					e.retryQueue[0].Quantity += 1
					e.retryQueue[0].RetrieveError = err
				} else {
					e.rmQueueErr()
				}
			}

		}
	}
}

func (e *Event) Publish(input Input) {
	e.bufferQueue = append(e.bufferQueue, input)
}
