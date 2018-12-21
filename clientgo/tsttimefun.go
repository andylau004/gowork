package main

import (
	"fmt"
	"sync"
	"time"
)

var g_wgsync sync.WaitGroup

func tst_time_1() {
	t1 := time.NewTimer(time.Second * 2)

	fmt.Println("curtime=", time.Now())
	<-t1.C
<<<<<<< HEAD
	fmt.Println("curtime=", time.Now(), " t1 is expired")
=======
	fmt.Println("t1 expire, curtime=", time.Now())
>>>>>>> ece6e0edf0e810f6b17cd2afddb84d50ca1dc763
}

func tst_time_2() {
	t2 := time.NewTimer(time.Second)
	// fmt.Println("curtime=", time.Now())

	g_wgsync.Add(1)
	go func() {
		fmt.Println("curtime=", time.Now(), " t2 is begin")
		<-t2.C
		fmt.Println("curtime=", time.Now(), " t2 is expired")
		// fmt.Println("t2 is expired")

		g_wgsync.Done()
	}()

	// stop2 := t2.Stop()
	// if stop2 {
	// 	fmt.Println("t2 is stopped")
	// }

	fmt.Println("wait beg")
	g_wgsync.Wait()
	fmt.Println("wait end")
	// fmt.Println("t1 expire, curtime=", time.Now())
}

<<<<<<< HEAD
func tst_time_3() {
	ticker := time.NewTicker(time.Second * 1)
	wrapGo(func() {

		for t := range ticker.C {
			fmt.Println("Ticker at: ", t)
		}

	})
	select {}
}
func entry_time_fun() {
	tst_time_1()
=======
func entry_time_fun() {
	tst_time_2()
>>>>>>> ece6e0edf0e810f6b17cd2afddb84d50ca1dc763
}
