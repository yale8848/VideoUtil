package main

import (
	"fmt"
	"os"
	"time"
)

func main()  {
	time.Sleep(1 * time.Second)
	if len(os.Args)<=2{
		fmt.Println("params error")
		return
	}
	arg1:=os.Args[1]
	arg2:=os.Args[2]
	fmt.Println("t.exe :  "+arg1+ " "+arg2)
}
