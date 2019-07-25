package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"sync"
	"syscall"
	"time"

	"./pool"
	"golang.org/x/sys/unix"
	// "reflect"
)

func GetCurTId() string {
	strId := fmt.Sprintf("%d ", unix.Gettid())
	return strId
}
func GetGoroutineID() uint64 {
	b := make([]byte, 64)
	runtime.Stack(b, false)
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)

	return n
}

func GetGoroutineIDStr() string {
	b := make([]byte, 64)
	runtime.Stack(b, false)
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)

	u10 := strconv.FormatUint(n, 10)
	// fmt.Printf("%T, %v\n", s10, s10)

	return fmt.Sprintf("%v", u10)
}

type T struct{}

// func (t *T) A() {
// 	fmt.Println( "aaaaaaaaaa" )
// }
// func (t *T) B() {
// 	fmt.Println( "bbbbbbbbbb" )
// }
func (t T) A() {
	fmt.Println("aaaaaaaaaa")
}
func (t T) B() {
	fmt.Println("bbbbbbbbbb")
}

type Ter interface {
	A()
	B()
}

func identity(z *T) *T {
	return z
}
func ref(z T) *T {
	return &z
}

func tst15() {

	var val int
	val = 2
	workcount := 0
	for i := 0; i < 15; i++ {
		val *= 2
		workcount++
	}

	fmt.Println("val=", val)
	fmt.Println("workcount=", workcount)

	var other int
	other = 2 << 4
	fmt.Println("other=", other)

	var tmp1 int
	tmp1 = -1
	tmp1 = (-1 << 3)
	fmt.Println("tmp1=", tmp1)
}

const srvAddr string = "172.17.0.2:444"

// connection pool
var g_connpoolObj pool.Pool

// < 连接对象, 使用次数 >
var g_mapUseCount sync.Map

func clientImpl(wg *sync.WaitGroup) {
	defer wg.Done()

	curId := unix.Gettid()

	v, err := g_connpoolObj.Get()
	if err != nil {
		fmt.Println(curId, " fatal error Get conn failed!!! err=", err)
		return
	}
	defer g_connpoolObj.Put(v)

	newConn := v.(net.Conn)
	fmt.Println(curId, " Get Conn Obj=", newConn)

	{
		newConn.Write([]byte("hello server\n"))
		// sendlen, err := newConn.Write([]byte("hello server\n"))
		// fmt.Println(curId, " sendlen=", sendlen, ", err=", err)
		respBuf := make([]byte, 2048)
		newConn.Read(respBuf)
		// recLen, err := newConn.Read(respBuf)
		// fmt.Println(curId, " recLen=", recLen, " respBuf=", string(respBuf))
	}

}

func clientFun() {
	var err error
	//factory 创建连接的方法
	factory := func() (interface{}, error) { return net.Dial("tcp", srvAddr) }

	//close 关闭连接的方法
	close := func(v interface{}) error { return v.(net.Conn).Close() }

	//创建一个连接池： 初始化5，最大连接30
	poolConfig := &pool.Config{
		InitialCap: 5,
		MaxCap:     30,
		Factory:    factory,
		Close:      close,
		//连接最大空闲时间，超过该时间的连接 将会关闭，可避免空闲时连接EOF，自动失效的问题
		IdleTimeout: 15 * time.Second,
	}

	g_connpoolObj, err = pool.NewChannelPool(poolConfig, &g_mapUseCount)
	if err != nil {
		fmt.Println("Create channelpool err=", err)
		return
	}
	lenpool := g_connpoolObj.Len()
	fmt.Println("conn pool len=", lenpool)

	tCount := 10
	wg := &sync.WaitGroup{}
	wg.Add(tCount)

	for i := 0; i < tCount; i++ {
		go clientImpl(wg)
	}
	wg.Wait()

	lenpool = g_connpoolObj.Len()
	fmt.Println("conn pool len=", lenpool)

	//释放连接池中的所有连接
	//p.Release()

	g_connpoolObj.UseCount()
}

func TstTcpConnPool() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGUSR1, syscall.SIGUSR2)

	clientFun()
	fmt.Println("use ctrl + c exit")
	<-c
	fmt.Println("all work done")
}

func tst1() {
	fmt.Println("tst1 .............")
}
func tst2() {
	fmt.Println("tst2 .............")
}

func TstDefer() {
	i := 1

	{
		if i == 1 {
			defer tst2()
			fmt.Println("abcdefd")
		} else {
			// defer tst2()
		}
	}
	{
		if i == 1 {
			defer tst1()
			fmt.Println("fghijklmn")
		} else {
			// defer tst2()
		}

	}
	fmt.Println("12345677")

}

type TransferFileBean struct {
	FileBuffer []byte
	FileSize   int
}

var g_chFileBean chan *TransferFileBean = make(chan *TransferFileBean, 1000)

// var g_chFileBean chan TransferFileBean = make(chan TransferFileBean)

func ConsumeData(wg *sync.WaitGroup) {

	// for {
	// 	select {
	// 	// case recvBean := <-g_chFileBean:
	// 	// 	fmt.Println("recvBean buffer=", &byteBuffer1, ", recvBean len=", recvBean.)
	// 	}
	// }

	recvBean := <-g_chFileBean
	fmt.Printf("ConsumeData buffer1=%p len1=%d buffer1=%s\n",
		&(recvBean.FileBuffer), recvBean.FileSize, string(recvBean.FileBuffer))
	wg.Done()

	{
		recvBean := <-g_chFileBean
		// fmt.Println("recvBean buffer=", &(recvBean.FileBuffer), ", recvBean FileSize=", recvBean.FileSize)
		fmt.Printf("ConsumeData recvBean2 buffer=%p len2=%d buffer2=%s\n",
			&(recvBean.FileBuffer), recvBean.FileSize, string(recvBean.FileBuffer))
		wg.Done()
	}

}
func TstChByte() {

	wg := &sync.WaitGroup{}
	wg.Add(2)

	{
		byteBuffer1 := make([]byte, 1024)
		copy(byteBuffer1, []byte("12345"))
		fmt.Printf(" byteBuffer1=%p\n", &(byteBuffer1))

		// var oneBean TransferFileBean
		oneBean := &TransferFileBean{}
		oneBean.FileBuffer = byteBuffer1
		oneBean.FileSize = len(string(oneBean.FileBuffer))
		len1 := oneBean.FileSize

		fmt.Printf(" oneBean1.FileBuffer=%p , len1=%d\n", &(oneBean.FileBuffer), len1)

		// fmt.Println("before1")
		g_chFileBean <- oneBean
		// fmt.Println("over1")
	}

	{
		byteBuffer2 := make([]byte, 1024)
		copy(byteBuffer2, []byte("678910"))
		// len2 := len(byteBuffer2)
		fmt.Printf(" byteBuffer2=%p\n", &(byteBuffer2))

		// var oneBean TransferFileBean
		oneBean := &TransferFileBean{}
		oneBean.FileBuffer = byteBuffer2[0:len(byteBuffer2)]
		// oneBean.FileSize = len2
		oneBean.FileSize = len(oneBean.FileBuffer)
		len2 := oneBean.FileSize

		fmt.Printf(" oneBean2.FileBuffer=%p , len2=%d\n", &(oneBean.FileBuffer), len2)

		// fmt.Println("before2")
		g_chFileBean <- oneBean
		// fmt.Println("after2")
	}

	go ConsumeData(wg)
	wg.Wait()
}

func TstThreadId() {

	workCount := 20
	wg := &sync.WaitGroup{}
	wg.Add(workCount)

	for i := 0; i < workCount; i++ {

		go func() {
			defer wg.Done()
			strCurId := GetGoroutineIDStr()

			for {
				time.Sleep(2 * time.Second)
				fmt.Println(strCurId, " is working")
			}

		}()
	}

	wg.Wait()
}

func main() {
	TstReflectEntry()
	return

	TstCtxEntry()
	return

	TstChanEntry()
	return

	// TstChanEntry()
	// return

	TstCtx()
	return

	TstThreadId()
	return

	TstTick()
	return

	TstCloseCh()
	return

	Tst_MRecver_MSender()
	return

	fmt.Println("GetGoroutineIDStr=", GetGoroutineIDStr())
	return
	// Wrrap()
	// return

	Tst_MRecver_MSender()
	return

	Tst1Recver_NSender()
	return

	TstChByte()
	return

	TstDefer()
	return
	// {
	// 	var mapUse sync.Map
	// 	mapUse.Store(1, 1)

	// 	//Load 方法，获得value
	// 	if v, ok := mapUse.Load(1); ok {
	// 		fmt.Println("v=", v)
	// 		mapUse.Store(1, v.(int)+1)

	// 		vt, _ := mapUse.Load(1)
	// 		fmt.Println("vt=", vt)
	// 	} else {

	// 	}
	// 	return
	// }
	TstTcpConnPool()
	return

	TstDinner()
	return

	tst15()
	return

	StartRecvUpload()
	return

	TstChanEntry()
	return

	TstBlg4Fun()

	time.Sleep(2 * time.Second)
	return

	tst_fun_entry()
	return

	var obj T
	p1 := &obj

	// p1.A()
	// p1.B()

	_ = *identity(p1)
	_ = *ref(obj)

}
