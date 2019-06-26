package main

import (
	"context"
	"fmt"
	"time"
	"unsafe"
)

type Human struct {
	name  string
	age   int
	phone string
}

type Student struct {
	Human
	school string
	fenshu int
}

type Empolyee struct {
	Human
	company string
	money   int
}

//Human实现SayHi方法
func (h Human) SayHi() {
	fmt.Printf("Hi, I am %s you can call me on %s\n", h.name, h.phone)
}

//Human实现Sing方法
func (h Human) Sing(lyrics string) {
	fmt.Println("La la la la...", lyrics)
}

//Employee重载Human的SayHi方法
func (e Empolyee) SayHi() {
	fmt.Printf("Hi, myname=%s my company=%s my phonenumber=%s\n", e.name, e.company, e.phone)
}

// Interface Men被Human,Student和Employee实现
// 因为这三个类型都实现了这两个方法
type Men interface {
	SayHi()
	Sing(lyrics string)
}

func timeoutHandler() {
	// 创建继承Background的子节点Context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	go doSth(ctx)

	//模拟程序运行 - Sleep 10秒
	time.Sleep(10 * time.Second)
	cancel() // 3秒后将提前取消 doSth goroutine
}

//每1秒work一下，同时会判断ctx是否被取消，如果是就退出
func doSth(ctx context.Context) {
	var i = 1
	for {
		time.Sleep(1 * time.Second)
		select {
		case <-ctx.Done():
			fmt.Println("done")
			return
		default:
			fmt.Printf("work %d seconds: \n", i)
		}
		i++
	}
}

func tst_Unsafe_point() {
	num := 5
	numPointer := &num
	flnum := (*float32)(unsafe.Pointer(numPointer))
	fmt.Println("flnum", flnum)
}

func tst_I_ex() {
	tst_Unsafe_point()
	return

	timeoutHandler()
	return

	mike := Student{Human{"Mike", 25, "222-222-XXX"}, "MIT", 1213}
	paul := Student{Human{"Paul", 26, "111-222-XXX"}, "Harvard", 1222}
	sam := Empolyee{Human{"Sam", 36, "444-222-XXX"}, "Golang Inc.", 33}
	Tom := Empolyee{Human{"Tom", 37, "222-444-XXX"}, "Things Ltd.", 444}

	//定义Men类型的变量i
	var i Men

	//i能存储Student
	i = mike
	fmt.Println("This is Mike, a Student:")
	i.SayHi()
	i.Sing("November rain")

	//i也能存储Employee
	i = Tom
	fmt.Println("This is Tom, an Employee:")
	i.SayHi()
	i.Sing("Born to be wild")

	//定义了slice Men
	fmt.Println("Let's use a slice of Men and see what happens")
	x := make([]Men, 3)
	//这三个都是不同类型的元素，但是他们实现了interface同一个接口
	x[0], x[1], x[2] = paul, sam, mike

	for _, value := range x {
		value.SayHi()
	}

}
