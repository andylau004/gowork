package main

import (
	// "context"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
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
			} // end select
		} // end for
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

// 测试
// 多个发送者，使用channel发送数据;
// 一个接收者，接受处理数据；
// 关闭channel的时候，由接收者关闭，也就是一个线程关闭，其他线程，发现关闭后，退出，
// 如果发送者关闭，多次关闭会导致，异常
func Tst1Recver_NSender() {
	rand.Seed(time.Now().UnixNano())
	// strId := GetCurTId()

	const MaxV = 1000
	const NumSenders = 5

	dataCh := make(chan int, 1000)
	stopCh := make(chan struct{})

	wg := &sync.WaitGroup{}
	wg.Add(NumSenders + 1)

	// sender 多个
	for i := 0; i < NumSenders; i++ {
		go func() {
			defer wg.Done()
			strId := GetCurTId()

			for {
				select {
				case <-stopCh:
					fmt.Println(strId, "sender thread exit ")
					return
				case dataCh <- rand.Intn(MaxV):
				}
			}
		}()
	}

	// receive 一个
	go func() {
		strId := GetCurTId()
		defer wg.Done()
		for {
			select {
			case v, _ := <-dataCh:
				fmt.Println("recver val=", v)
				if v == MaxV-1 {
					close(stopCh)
					fmt.Println(strId, "recver thread exit")
					// wg.Done()
					return
				}
			}
		}
	}()

	wg.Wait()
	fmt.Println("all work done")
}

// 测试多sender，多receiver
func Tst_MRecver_MSender() {

	rand.Seed(time.Now().UnixNano())

	const MaxV = 1000
	const NumSenders = 50
	const NumReceivers = 10

	dataCh := make(chan int, 1000)
	stopCh := make(chan struct{})

	toStop := make(chan string, NumSenders+NumReceivers)

	wg := &sync.WaitGroup{}
	wg.Add(NumSenders + NumReceivers)

	var stoppedBy string
	go func() {
		stoppedBy = <-toStop
		close(stopCh)
	}()

	// sender 多个
	for i := 0; i < NumSenders; i++ {

		go func(id string) {

			defer wg.Done()
			strId := GetGoroutineIDStr()

			for {
				value := rand.Intn(MaxV)
				if value == 0 {
					toStop <- "sender ### " + id
					return
				}

				select {
				case <-stopCh:
					fmt.Println(strId, "sender thread exit ")
					return
				case dataCh <- value:
				}
			}

		}(strconv.Itoa(i))

	}

	// receive 多个
	for i := 0; i < NumReceivers; i++ {

		go func(id string) {
			strId := GetGoroutineIDStr()
			defer wg.Done()

			for {
				select {
				case <-stopCh:
					return
				case v := <-dataCh:
					// fmt.Println("recver val=", v)
					if v == MaxV-1 {
						select {
						case toStop <- "receiver#" + id:
							{
								fmt.Println(strId, "Receiver thread exit ")
							}
						default:
						}
						return
					}

				}

			}
		}(strconv.Itoa(i))

	}

	wg.Wait()
	fmt.Println("all work done")

}

func Wrrap() {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer func() {
			fmt.Println("exit fun")
			wg.Done()
		}()
		i := 1
		for {
			i++
			if i == 100 {
				fmt.Println("i==100, i=", i)
				return
			}
		}
	}()

	wg.Wait()
	fmt.Println("all work done")
}

func TstCloseCh() {
	ch := make(chan int, 5)
	ch <- 18
	close(ch)

	v, ok := <-ch
	if ok {
		fmt.Println("v=", v)
	}

	{
		v, ok = <-ch
		if !ok {
			fmt.Println("channle closed, can't read")
		}
	}

}

func TickWork() {
	ticker := time.Tick(2 * time.Second)

	for {

		select {
		case <-ticker:
			fmt.Println("execute 1s work")
			// default:
			// 	fmt.Println("execute default")
		}

	}

}
func TstTick() {
	chExit := make(chan int)
	go TickWork()
	<-chExit
}

func TstSelect() {
	ch := make(chan int)

	fmt.Println("before select")
	select {
	case tmpV := <-ch:
		fmt.Println("tmpV=", tmpV)
		// default:
		// 	fmt.Println("default  exec ")
	}
	fmt.Println("TstSelect done")
	return
}

func TstChanClose() {
	chExit := make(chan bool)
	wg := &sync.WaitGroup{}
	wg.Add(5)

	for i := 0; i < 5; i++ {

		go func() {
			for {
				select {
				case <-chExit:
					fmt.Println("tId=", GetGoroutineIDStr(), " exit")
					wg.Done()
					return
				default:
					fmt.Println("tId=", GetGoroutineIDStr(), " working...")
				}
				time.Sleep(time.Millisecond * 300)
			}
		}()

	}

	time.Sleep(time.Second * 2)
	fmt.Println("sleep compelte")

	for i := 0; i < 5; i++ {
		fmt.Println("before i=", i)
		chExit <- true
		fmt.Println("after i=", i)
	}

	wg.Wait()
	// time.Sleep(time.Second * 2)
}
func TstChanEntry() {
	TstChanClose()
	return

	TstSelect()
	return

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
