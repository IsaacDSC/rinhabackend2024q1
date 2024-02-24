package tests

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"rinhabackend/external/queue"
	"testing"
	"time"
)

type Account struct {
	ID string
}

type ISync struct {
}

func (i ISync) Consumer(input []byte) error {
	var un Account
	err := json.Unmarshal(input, &un)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", un)
	return errors.New("TESTES")
}

func (i ISync) ConsumerErr(input queue.Retry) error {
	fmt.Printf("WithErr: %+v\n", input)
	return errors.New("TESTES")
}

func TestQueue(t *testing.T) {
	account := Account{ID: uuid.NewString()}
	input, _ := json.Marshal(account)
	isyncq := queue.NewEvent(new(ISync))
	go isyncq.Consume()
	isyncq.Publish(input)
	time.Sleep(time.Second * 14)
}
