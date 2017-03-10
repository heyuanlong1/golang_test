package main

import "fmt"

func main() {
	cs := make([](chan int), 10)
	for i := 0; i < len(cs); i++ {
		cs[i] = make(chan int)
	}

	for i := range cs {
		go func() {
			cs[i] <- i
		}()
	}
	/*
	   for i := range cs {
	       go func(index int) {
	           cs[index] <- index
	       }(i)
	   }
	*/

	for i := 0; i < len(cs); i++ {
		t := <-cs[i]
		fmt.Println(t)
	}
}
