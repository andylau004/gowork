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

type ServeHandler func(c net.Conn) error

type WorkerPool struct {
	// Function for serving server connections.
	// It must leave c unclosed.
	WorkerFunc ServeHandler

	lock         sync.Mutex
	workersCount int
	mustStop     bool

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
	fmt.Println("new add client, addr=", c.RemoteAddr().String())
	return true
}

func (wp *WorkerPool) getCh() *workerChan {
	var ch *workerChan
	// createWorker := false

	wp.lock.Lock()
	ready := wp.ready
	n := len(ready) - 1
	if n < 0 {
		fmt.Println("111 n=", n)
		// if wp.workersCount < wp.MaxWorkersCount {
		// createWorker = true
		wp.workersCount++
		// }
	} else {
		fmt.Println("222 n=", n)
		ch = ready[n]
		ready[n] = nil
		wp.ready = ready[:n]
	}
	wp.lock.Unlock()

	if ch == nil {
		// if !createWorker {
		// 	return nil
		// }
		fmt.Println("ch == nil")
		vch := wp.workerChanPool.Get()
		if vch == nil {
			vch = &workerChan{
				ch: make(chan net.Conn, 11),
			}
		}
		fmt.Println("vch=", vch)
		ch = vch.(*workerChan)
		fmt.Printf("ch=%+v\n", *ch)
		go func() {
			wp.workerFunc(ch)
			wp.workerChanPool.Put(vch)
		}()
	}
	return ch
}

func (wp *WorkerPool) workerFunc(ch *workerChan) {
	defer fmt.Printf(" worker func end, ch=%+v\n", *ch)

	var c net.Conn
	var err error
	for c = range ch.ch {
		if c == nil {
			fmt.Println("c == nil! err=", err)
			break
		}

		if err = wp.WorkerFunc(c); err != nil {
			fmt.Println(" WorkFunc failed! err=", err)
		}
		c = nil

		if !wp.release(ch) {
			fmt.Println("wp release failed! err=", err)
			break
		}
	}
	wp.lock.Lock()
	wp.workersCount--
	wp.lock.Unlock()
}

func (wp *WorkerPool) release(ch *workerChan) bool {
	ch.lastUseTime = time.Now()

	wp.lock.Lock()
	defer wp.lock.Unlock()
	if wp.mustStop {
		return false
	}
	fmt.Println("before cur ready len = ", len(wp.ready))
	wp.ready = append(wp.ready, ch)
	fmt.Println("after cur ready len = ", len(wp.ready))
	// wp.lock.Unlock()
	return true
}

func ConnHandler(conn net.Conn) error {
	fmt.Println("conn handler beg")
	defer fmt.Println("conn handler end")

	defer conn.Close()

	var buf [2048]byte
	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			fmt.Println("read client failed! err=", err)
			return err
		}
		fmt.Println("read from client=", conn.RemoteAddr().String, ", mgs=", string(buf[0:n]))
		_, err2 := conn.Write(buf[0:n])
		if err2 != nil {
			fmt.Println("write client failed! err=", err)
			return err
		}
	}

	return nil
}

func WorkerPoolEntry() {

	wp := &WorkerPool{
		WorkerFunc: ConnHandler,
		// MaxWorkersCount: 1000,
		// LogAllErrors:    s.LogAllErrors,
		// Logger:          s.logger(),
		// connState:       s.setState,
	}
	wp.Start()

	addr := ":9981"
	ln, err := net.Listen("tcp4", addr)
	if err != nil {
		fmt.Println("listen failed! err=", err)
		return
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil { // handle error
			fmt.Println("accept failed! err=", err)
			return
		}
		wp.Serve(conn)
	}

}
