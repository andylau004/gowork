package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

func dummy_workerpool() {

	fmt.Println("")
}

type WorkerPool struct {
	lock           sync.Mutex
	workersCount   int
	mustStop       bool
	ready          []*workerChan
	stopCh         chan struct{}
	workerChanPool sync.Pool
}

type workerChan struct {
	lastUseTime time.Time
	ch          chan net.Conn
}

func (wp *WorkerPool) Start() {
	if wp.stopCh != nil {
		panic("BUG: workerPool already started")
	}
	wp.stopCh = make(chan struct{})
	stopCh := wp.stopCh
	fmt.Printf("wp.stopCh=%+v\n", wp.stopCh)
	fmt.Printf("stopCh=%+v\n", stopCh)
}

func (wp *WorkerPool) Serve(c net.Conn) bool {
	ch := wp.getCh()
	if ch == nil {
		fmt.Println("get ch == nil")
		return false
	}
	ch.ch <- c
	return true
}

func (wp *WorkerPool) getCh() *workerChan {
	var ch *workerChan
	// createWorker := false

	wp.lock.Lock()
	ready := wp.ready
	n := len(ready) - 1
	if n < 0 {
		// if wp.workersCount < wp.MaxWorkersCount {
		// createWorker = true
		wp.workersCount++
		// }
	} else {
		ch = ready[n]
		ready[n] = nil
		wp.ready = ready[:n]
	}
	wp.lock.Unlock()

	if ch == nil {
		// if !createWorker {
		// 	return nil
		// }
		vch := wp.workerChanPool.Get()
		if vch == nil {
			vch = &workerChan{
				ch: make(chan net.Conn, 10),
			}
		}
		ch = vch.(*workerChan)
		go func() {
			wp.workerFunc(ch)
			wp.workerChanPool.Put(vch)
		}()
	}
	return ch
}

func (wp *WorkerPool) workerFunc(ch *workerChan) {
	var c net.Conn

	var err error
	for c = range ch.ch {
		if c == nil {
			break
		}

	}

}

func WorkerPoolEntry() {

}
