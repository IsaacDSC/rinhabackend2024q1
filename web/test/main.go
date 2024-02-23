package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	fmt.Println("wait to cancel", ctxTimeout)
	select {
	case <-time.After(5 * time.Second):
		fmt.Println("Slow process completed")
	case <-ctxTimeout.Done():
		fmt.Println("Process timed out")
	}
	//time.Sleep(time.Minute * 6)
}
