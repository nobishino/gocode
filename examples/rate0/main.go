package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/time/rate"
)

func main() {
	r := rate.Every(1 * time.Second)

	limiter := rate.NewLimiter(r, 1)

	for {
		limiter.Wait(context.Background())
		f()
	}
}

func f() {
	fmt.Println("DO")
}
