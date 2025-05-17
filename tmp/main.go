package main

import (
	"context"
	"fmt"
	"time"
)

func tickOperation(ctx context.Context) {
	ticker := time.NewTicker(200 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			fmt.Println("tickOperation: ", time.Now().UnixNano())
		case <-ctx.Done():
			fmt.Println("Cancel: ", time.Now())
			return
		}
	}
}

func main() {
	ctx, cancel2 := context.WithCancel(context.Background())
	go tickOperation(ctx)

	time.Sleep(5 * time.Second)
	cancel2()
	time.Sleep(1 * time.Second)
}

func testValueContext() {
	type key int
	const EmailKey key = 0
	ctx := context.Background()
	ctxWithValue := context.WithValue(ctx, EmailKey, "asdqwe@aasd.ru")

	if userEmail, ok := ctxWithValue.Value(EmailKey).(string); !ok {
		fmt.Println("Email key is not a string")
	} else {
		fmt.Println("User email: " + userEmail)
	}
}

func testTimoutContext() {
	ctx := context.Background()
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	done := make(chan struct{})
	go func() {
		time.Sleep(3 * time.Second)
		close(done)
	}()

	select {
	case <-done:
		fmt.Println("Done task")
	case <-ctxWithTimeout.Done():
		fmt.Println("Timeout")
	}
}
