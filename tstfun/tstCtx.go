package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"time"

	"net/http"
	_ "net/http/pprof"

	blg4go "github.com/YoungPioneers/blog4go"
)

func dummyFunc_ctx() {
	fmt.Println()
	blg4go.Info("")
}

func C(ctx context.Context) string {
	select {
	case <-ctx.Done():
		return "C Done"
	}
	return ""
}
func B(ctx context.Context) string {
	ctx, _ = context.WithCancel(ctx)
	go log.Println(C(ctx))
	select {
	case <-ctx.Done():
		return "B Done"
	}
	return ""
}

func A(ctx context.Context) string {
	go log.Println(B(ctx))
	select {
	case <-ctx.Done():
		return "A Done"
	}
	return ""
}

// 模拟一个最小执行时间的阻塞函数
func inc(a int) int {
	res := a + 1                // 虽然我只做了一次简单的 +1 的运算,
	time.Sleep(1 * time.Second) // 但是由于我的机器指令集中没有这条指令,
	// 所以在我执行了 1000000000 条机器指令, 续了 1s 之后, 我才终于得到结果。B)
	return res
}

func Add(ctx context.Context, a, b int) int {
	ret := 0

	for i := 0; i < a; i++ {
		ret = inc(ret)
		select {
		case <-ctx.Done():
			return -1
		default:
		}
	}

	for i := 0; i < b; i++ {
		ret = inc(ret)
		select {
		case <-ctx.Done():
			return -1
		default:
		}
	}
	return ret
}

func TstCtx() {
	go http.ListenAndServe(":8989", nil)

	// 使用开放的 API 计算 a+b
	// a := 1
	// b := 2
	// timeout := 6 * time.Second
	// ctx, _ := context.WithTimeout(context.Background(), timeout)
	// res := Add(ctx, 1, 2)
	// fmt.Printf("Compute: %d+%d, result: %d\n", a, b, res)

	{
		// 手动取消
		a := 1
		b := 2
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			time.Sleep(2 * time.Second)
			fmt.Println("before cancel")
			cancel() // 在调用处主动取消
			fmt.Println("after cancel")
		}()
		res := Add(ctx, 1, 2)
		fmt.Printf("Compute: %d+%d, result: %d\n", a, b, res)
	}
}

func chiHanbao(ctx context.Context) <-chan int {
	c := make(chan int)

	n := 0
	t := 0

	go func() {

		for {
			select {
			case <-ctx.Done():
				fmt.Printf("耗时 %d 秒，吃了 %d 个汉堡 \n", t, n)
				return
			case c <- n:
				incr := rand.Intn(5)
				n += incr
				if n >= 10 {
					n = 10
				}
				t++
				fmt.Printf("eat %d hanbaoge\n", n)
			}
		}

	}()

	return c
}

func TstCtx_WithCancel() {
	ctx, cancel := context.WithCancel(context.Background())

	eatNum := chiHanbao(ctx)
	for n := range eatNum {
		if n >= 10 {
			cancel()
			break
		}
	}

	fmt.Println("正在统计结果。。。")
	time.Sleep(1 * time.Second)
}

func chiHanbao_WithTimeOut(ctx context.Context) {
	n := 0
	for {
		select {
		case <-ctx.Done():
			fmt.Println("ctx Done")
			return
		default:
			incr := rand.Intn(5)
			n += incr
			fmt.Printf("eat %d ge hanbao\n", n)
		}
		time.Sleep(time.Second)
	}
}
func TstCtx_WithTimeout() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	chiHanbao_WithTimeOut(ctx)
	defer cancel()
}

func process(ctx context.Context) {
	session, ok := ctx.Value("session").(int)
	fmt.Println("ok1=", ok)

	if !ok {
		fmt.Println("something wrong")
		return
	}
	if session != 1 {
		fmt.Println("session 未通过")
		return
	}
	traceID := ctx.Value("trace_id").(string)
	fmt.Println("traceID:", traceID, " session:", session)
}
func TstCtx_WithValue() {
	ctx := context.WithValue(context.Background(), "trace_id", "88888888")

	// 携带session到后面的程序中去
	ctx = context.WithValue(ctx, "session", 1)
	process(ctx)
}

func watch(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println(name, "监控退出，停止了...")
			return
		default:
			fmt.Println(name, "goroutine监控中...")
			time.Sleep(2 * time.Second)
		}
	}
}

func Tst_Cancel_1() {
	ctx, cancel := context.WithCancel(context.Background())
	go watch(ctx, "【监控1】")
	go watch(ctx, "【监控2】")
	go watch(ctx, "【监控3】")
	time.Sleep(10 * time.Second)
	fmt.Println("可以了，通知监控停止")
	cancel()
	time.Sleep(5 * time.Second)
}

func Tst_Ctx_WithTimOut_1() {
	// ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	done := make(chan int, 1)
	go func() {
		// time.Sleep(time.Second)
		time.Sleep(200 * time.Millisecond)
		done <- 1
	}()

	select {
	case <-done:
		fmt.Println("channel done ont time")
	case <-ctx.Done():
		fmt.Println("ctx timeout, ctx Err=", ctx.Err())
	}

	fmt.Println("main thread Exit...")
}

func mainTask(ctx context.Context, taskName string) {

}

func Tst_Ctx_12() {
	ctx, _ := context.WithTimeout(context.Background(), 50*time.Millisecond)

	go mainTask(ctx, "1")
	go mainTask(ctx, "2")
	go mainTask(ctx, "3")

	select {
	case <-ctx.Done():
		fmt.Println("main error:", ctx.Err())
	}

	fmt.Println("main exit...")
	time.Sleep(3 * time.Second)
}

type Result struct {
	r   *http.Response
	err error
}

func tstGetBaidu() {
	// ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tr := &http.Transport{}
	client := &http.Client{Transport: tr}

	// resultChan := make(chan Result, 1)
	resultChan := make(chan Result)
	req, err := http.NewRequest("GET", "http://www.baidu.com", nil)
	// req, err := http.NewRequest("GET", "http://www.google.com", nil)
	if err != nil {
		fmt.Println("http request failed , err=", err)
		return
	}

	go func() {
		tStart := time.Now()
		resp, err := client.Do(req)
		fmt.Println("client.Do timeinterval=", time.Now().Sub(tStart))

		pack := Result{r: resp, err: err}
		resultChan <- pack

		fmt.Println("cheRes <----- pack done")
		fmt.Printf("pack=%+v\n", pack)
	}()

	select {
	case <-ctx.Done():
		tr.CancelRequest(req)
		er := <-resultChan
		fmt.Printf("\n\n")
		fmt.Printf("Timeout!!! er=%+v\n", er)

	case res := <-resultChan:
		defer res.r.Body.Close()
		out, _ := ioutil.ReadAll(res.r.Body)
		fmt.Println("Server Response len=", len(string(out)), ", respBody=", string(out))
	}
	fmt.Println("main thread exit, ctx=", ctx)
}

func TstCtxEntry() {
	tstGetBaidu()
	return

	Tst_Ctx_WithTimOut_1()
	return

	Modify_Map()
	return

	Modify_Addr()
	return

	Tst_Cancel_1()
	return

	TstCtx_WithValue()
	return

	TstCtx_WithTimeout()
	return

	TstCtx_WithCancel()
	return

}
