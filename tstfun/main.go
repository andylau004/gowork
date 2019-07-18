package main

import (
	"fmt"
	"time"
	// "reflect"
)

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

// 检查堆栈逃逸
func main() {

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
