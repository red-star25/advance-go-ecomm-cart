package main

import "fmt"

func main() {
	var myChan = make(chan int)

	go func() {
		myChan <- 1
	}()

	val := <-myChan

	fmt.Println(val)
}
