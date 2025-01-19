package main

import "fmt"

/*
Channel: Used to acheive concurrency in Go.
- They communicate and transfer data between goroutines.
- The communication are bidectional meaning you can send and receive data from the same go-routine.
*/
func main() {
	var myChan = make(chan int)

	go func() {
		myChan <- 1
	}()

	val := <-myChan

	fmt.Println(val)
}
