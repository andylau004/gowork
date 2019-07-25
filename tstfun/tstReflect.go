package main

import (
	"fmt"
	"reflect"

	blg4go "github.com/YoungPioneers/blog4go"
)

func dummyFunc_Reflect() {
	fmt.Println()
	blg4go.Info("")
}

func tstReflect_1() {

	a := 1
	v := reflect.ValueOf(a)
	fmt.Println("v type=", v.Type())
	fmt.Println("v type=", v.CanSet())

	v = reflect.ValueOf(&a)
	fmt.Println("v Type:", v.Type())
	fmt.Println("v CanSet:", v.CanSet())

	v = v.Elem() // element value
	fmt.Println("v Type:", v.Type())
	fmt.Println("v CanSet:", v.CanSet())

	// set
	v.SetInt(2)
	fmt.Println("after set, v:", v)

	newValue := reflect.ValueOf(3)
	v.Set(newValue)
	fmt.Println("after set, v:", v)

}

func tstReflectSlice() {
	a := []int{1, 2}
	fmt.Println("a=", a)
}

type Orange struct {
	Size int
	Name string
}

func tstReflectStruct() {
	oneBean := &Orange{
		Size: 123,
		Name: "tstName",
	}
	v := reflect.ValueOf(oneBean)
	fmt.Println("v=", v)
}

func TstReflectEntry() {
	tstReflectStruct()
	return

	tstReflectSlice()
	return

	tstReflect_1()
	return

}
