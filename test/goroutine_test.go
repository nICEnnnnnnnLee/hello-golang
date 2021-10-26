package test

import (
	"fmt"
	"testing"
	"time"
)

func TestGoroutine(t *testing.T) {
	fmt.Println("main start")
	sig := make(chan int)

	go func() {
		fmt.Println("routine-1 start")
		defer fmt.Println("routine-1 end")
		go func() {
			fmt.Println("routine-2 start")
			defer fmt.Println("routine-2 end")
			count := 1
			for {
				time.Sleep(time.Second)
				fmt.Println("routine-2 sleep awake", count)
				count++
				if count > 20 {
					sig <- 1
				} else if count > 30 {
					break
				}
			}

		}()

	}()
	<-sig
	fmt.Println("main end")
}
