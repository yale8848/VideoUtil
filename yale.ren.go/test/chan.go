package main

import (
	"fmt"
	"time"
)

var ch = make(chan int)

func test()  {

	time.Sleep(3*time.Second)

	ch<-1
	v :=<- ch
	fmt.Println(v)

	ch<-2
}

func main()  {
	//ch<-2
	go test()

	fmt.Println("aa")

	v :=<- ch
	fmt.Println(v)


	fmt.Println("bb")
}

