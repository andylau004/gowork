package main

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/xxjwxc/gowp/workpool"
)

func dummy_wp() {
	fmt.Println("")
	time.Sleep(time.Second)
}

func WorkImpl(ii int) error {
	for j := 0; j < 10; j++ { // 0-10 values per print
		fmt.Println(fmt.Sprintf("%v->\t%v", ii, j))
		if ii == 1 {
			return errors.Cause(errors.New("my test err")) // have err return
		}
		// time.Sleep(1 * time.Second)
	}
	return nil
}

// type TaskHandler func(ii int) error

func Work1(ii int) error {
	for j := 0; j < 10; j++ {
		fmt.Println(fmt.Sprintf("%v->\t%v", ii, j))
		time.Sleep(1 * time.Millisecond)
	}
	// time.Sleep(1 * time.Second)
	return nil
}

func TstWorkPoolFun() {
	{
		wp := workpool.New(5) // Set the maximum number of threads
		for i := 0; i < 10; i++ {
			//	ii := i
			wp.Do(func() error {
				for j := 0; j < 5; j++ {
					// fmt.Println(fmt.Sprintf("%v->\t%v", ii, j))
					time.Sleep(1 * time.Millisecond)
				}
				return nil
			})

			fmt.Println(wp.IsDone())
		}
		wp.Wait()
		fmt.Println(wp.IsDone())
		fmt.Println("down")
		return
	}
	wp := workpool.New(10)

	for i := 0; i < 10; i++ { // Open 20 requests
		// ii := i
		// var p1 TaskHandler
		// p1 = Work1
		// pfnFun := p1(ii)
		// wp.Do(Work1(ii))
	} // for 0 - 10

	err := wp.Wait()
	if err != nil {
		fmt.Println("Wait found err!, err=", err)
	}
	fmt.Println("work done ...")
}
