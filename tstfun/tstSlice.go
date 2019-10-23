package main

import (
	"fmt"

	blg4go "github.com/YoungPioneers/blog4go"
)

func dummyFunc_slice() {
	fmt.Println()
	blg4go.Info("")
}

func tstSlice1() {

	arr := []string{"111", "222", "333"}
	for i, v := range arr {
		fmt.Println("i=", i, ", v=", v)
	}
}

type student struct {
	Name string
	Age  int
}

func ParseStudent() {
	m := make(map[string]*student)
	stus := []student{
		{Name: "zhang", Age: 22},
		{Name: "li", Age: 23},
		{Name: "wang", Age: 24},
	}
	for _, stu := range stus {
		fmt.Printf("&stu=%p\n", &stu)
		fmt.Println("stus=", stu)

		// m[stu.Name] = &stu
		newStu := stu
		m[stu.Name] = &newStu
	}
	for i, v := range m {
		fmt.Println("i=", i, ", v=", v)
	}
}

func tstSlice2() {
	ParseStudent()
}

func TstSliceEntry() {
	tstSlice2()
	return

	tstSlice1()
	return
}
