package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	query(ctx)

	time.Sleep(time.Hour)
}

func query(ctx context.Context) {
	notifyChan := make(chan string)
	doing := func() {
		time.Sleep(time.Second)
		fmt.Println("1:", time.Now())
		time.Sleep(time.Second * 3)
		fmt.Println("2:", time.Now())
		time.Sleep(time.Second * 5)
		fmt.Println("3:", time.Now())
		notifyChan <- "1"
	}
	go doing()
	select {
	case <-notifyChan:
		fmt.Println("done")
	case <-ctx.Done():
		fmt.Println("timeout")
	}
}
