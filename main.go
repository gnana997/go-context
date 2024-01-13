package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

func main() {

	ctx := context.WithValue(context.Background(), "deadline", 100)

	val, err := fetchThirdPartyData(ctx, 1)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Value is %d\n", val)

}

type Response struct {
	val int
	err error
}

func fetchThirdPartyData(ctx context.Context, id int) (int, error) {
	val := ctx.Value("deadline")
	duration := time.Duration(val.(int)) * time.Millisecond
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(duration))
	defer cancel()

	respch := make(chan Response)

	go func() {
		val, err := callThridPartyApi(id)
		respch <- Response{val, err}
	}()

	select {
	case <-ctx.Done():
		return 0, fmt.Errorf("thrid apI is taking more than 100ms")
	case resp := <-respch:
		return resp.val, resp.err
	}
}

func callThridPartyApi(id int) (int, error) {
	time.Sleep(500 * time.Millisecond)

	return id, nil
}
