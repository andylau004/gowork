package main

import (
	// "context"
	"fmt"
	"time"

	blg4go "github.com/YoungPioneers/blog4go"
)

func dummyFunc1() {
	fmt.Println()
	blg4go.Info("")
}

type Person struct {
	Name    string
	Age     uint8
	Address Addr
}
type Addr struct {
	city     string
	district string
}

// 测试channel传输复杂的Struct数据
func testTransferStruct() {
	personChan := make(chan Person, 1)

	onePerson := Person{"xiaoming", 10, Addr{"haidianqu", "zhaodenglu"}}
	personChan <- onePerson

	onePerson.Address = Addr{"xichengqu", "yingxionglu"}
	fmt.Printf("modifyPerson : %+v \n", onePerson)

	{
		recvPerson := <-personChan
		fmt.Printf("recvPerson : %+v \n", recvPerson)
	}

	// fy@fy:~/WorkDir/mygitwork/gowork/tstfun$ go build && ./tstfun
	// modifyPerson : {Name:xiaoming Age:10 Address:{city:xichengqu district:yingxionglu}}
	// recvPerson : {Name:xiaoming Age:10 Address:{city:haidianqu district:zhaodenglu}}
}

// 将多个输入的channel进行合并成一个channel
func TstMergeChan() {
	input1 := make(chan int)
	input2 := make(chan int)
	output := make(chan int)

	go func(in1, in2 <-chan int, out chan<- int) {

		for {
			select {
			case v := <-in1:
				out <- v
			case v := <-in2:
				out <- v
			} // end select
		} // end for

	}(input1, input2, output)

	go func() {
		for i := 0; i < 10; i++ {
			input1 <- i
			time.Sleep(time.Millisecond * 100)
		}
	}()
	go func() {
		for i := 20; i < 30; i++ {
			input2 <- i
			time.Sleep(time.Millisecond * 100)
		}
	}()

	go func() {
		for {
			select {
			case value := <-output:
				fmt.Println("输出：", value)
			}
		}
	}()
	time.Sleep(time.Second * 5)
	fmt.Println("主线程退出")

}

func TstChanExit() {
	ch := make(chan int)

	go func() {
		v := <-ch
		fmt.Println("v1111=", v)
	}()

	ch <- 1
	fmt.Println("2")
}

// 检查channel读写超时，并做超时的处理
func TstTimeOutChan() {
	g := make(chan int)
	quit := make(chan bool)

	go func() {
		for {
			select {
			case v := <-g:
				fmt.Println("val=", v)
			case <-time.After(time.Second * time.Duration(3)):
				fmt.Println("after logic start")
				quit <- true
				fmt.Println("work timeout ")
				return
			}

		}
	}()

	fmt.Println("aaaaaaaaaaaaaaaaaaaaa")
	for i := 0; i < 3; i++ {
		fmt.Println("beg i=", i)
		g <- i
		fmt.Println("end i=", i)
	}
	fmt.Println("bbbbbbbbbbbbbbbbbbbb")

	<-quit
	fmt.Println("all wokr done")
}

// 指定channel是输入还是输出型的，防止编写时写错误输入输出，指定了的话，可以在编译时期作错误的检查
func TstInAndOutChan() {
	ch := make(chan int)
	quit := make(chan bool)
	fmt.Println("out &ch=", &ch)

	go func(inChan chan<- int) {
		fmt.Println("in &ch=", &inChan)

		for i := 0; i < 10; i++ {
			inChan <- i
			// ch <- i
			time.Sleep(time.Millisecond * 100)
		}
		quit <- true
	}(ch)

	go func(outChan <-chan int) {
		fmt.Println("outChan &ch=", &outChan)
		for {
			select {
			case v := <-outChan:
				fmt.Println("print out value : ", v)
			case <-quit:
				fmt.Println("收到退出通知，退出")
				return

			} // end select
		} // end for
	}(ch)

	<-quit
	fmt.Println("all wokr done")
}

//测试通过channel来控制最大并发数，来处理事件
func TstMaxNumControl() {
	// maxNum := 3
	limit := make(chan bool, 3)
	quit := make(chan bool)

	for i := 0; i < 10; i++ {
		// fmt.Println("start worker : ", i)
		limit <- true

		go func(i int) {
			fmt.Println("do worker start: ", i)
			time.Sleep(time.Millisecond * 20)
			defer fmt.Println("do worker finish: ", i)

			<-limit

			if i == 9 {
				fmt.Println("完成任务")
				quit <- true
			}
		}(i)
	}

	<-quit
	fmt.Println("收到退出通知，主程序退出")
}

func TstChanEntry() {
	TstMaxNumControl()
	return

	TstInAndOutChan()
	return

	TstTimeOutChan()
	return

	TstChanExit()
	return

	TstMergeChan()
	return

	testTransferStruct()

}
