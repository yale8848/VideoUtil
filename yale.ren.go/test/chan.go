package main

import (
	"fmt"
	"time"
	"strconv"
)

var ch = make(chan int)

func test()  {
	for i:=0;;i++{
		time.Sleep(3*time.Second)
		fmt.Println("send "+strconv.Itoa(i))
		ch<-i
	}
}
func test1()  {
	time.Sleep(3*time.Second)
	//fmt.Println(<-ch)
}
func main()  {
	go test1()
	ch<-2


}
func main1()  {
	go test()

    for {
		v :=<- ch
		fmt.Println("receive "+strconv.Itoa(v))
	}


	fmt.Println("bb")
}

