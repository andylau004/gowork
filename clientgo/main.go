package main

import (
	"fmt"
	"net"
	"sync"
	"time"

	"./util"
)

var g_clientCount int
var g_wgWork sync.WaitGroup

var g_srvAddr = "218.241.161.57:10669"

func connectNgrokSrv(wgWrk *sync.WaitGroup) {
	defer wgWrk.Done()

	var err error
	var rawConn net.Conn

	if rawConn, err = net.Dial("tcp", g_srvAddr); err != nil {
		fmt.Printf("connect srv: %s failed!\n", g_srvAddr)
		return
	} else {
		fmt.Printf("connect srv: %s successed!\n", g_srvAddr)
	}
	defer rawConn.Close()

	select {}

}

func tst_tcp_client() {
	g_clientCount = 500
	g_wgWork.Add(g_clientCount)

	// fmt.Println("hello serv")

	for icount := 0; icount < g_clientCount; icount++ {
		go connectNgrokSrv(&g_wgWork)
	}

	g_wgWork.Wait()

}

var g_wgT1 sync.WaitGroup

func wrapGo(funcWork func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				err := util.MakePanicTrace(r)
				fmt.Println("crashed! err=", err)
			}
		}()
		funcWork()
	}()
}

func tst_1() {
	defer func() {
		if r := recover(); r != nil {
			err := util.MakePanicTrace(r)
			fmt.Println("crashed! err=", err)
		}
	}()

	objCh := make(chan int)
	defer close(objCh)

	fmt.Println("aaaaa\n")
	close(objCh)
	go func() {
		objCh <- (5 + 2)
	}()

	outVal := <-objCh
	fmt.Println("bbbbbb, outval=", outVal)
}

func tst_2() {

	g_wgT1.Add(1)
	chtmp := make(chan int)

	wrapGo(func() {
		fmt.Println("push it")
		chtmp <- 1
		fmt.Println("push done")
	})

	wrapGo(func() {
		fmt.Println("aaaaaaaaaaaaaaaaa")

		fmt.Printf("firstval=%d\n", <-chtmp)
		// close(chtmp)
		chtmp <- 2
		fmt.Printf("secondval=%d\n", <-chtmp)

		time.Sleep(2 * time.Second)
		g_wgT1.Done()
	})

	g_wgT1.Wait()
	fmt.Println("ccccccccccccc")
}

func main() {

	entry_time_fun()

	return

	tst_2()
	// tst_1()

	fmt.Println("game over")
}
