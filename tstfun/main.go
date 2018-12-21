package main

import (
	"fmt"
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

// 检查堆栈逃逸
func main() {

	var obj T
	p1 := &obj

	// p1.A()
	// p1.B()

	_ = *identity(p1)
	_ = *ref(obj)

}
