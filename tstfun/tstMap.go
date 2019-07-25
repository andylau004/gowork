package main

import (
	"fmt"
)

func modify(ip *int) {
	fmt.Printf("函数里接收到的指针的内存地址是：%p\n", ip)
	*ip = 1
}
func Modify_Addr() {
	i := 10
	ip := &i
	fmt.Printf("原始指针的内存地址是：%p\n", ip)
	modify(ip)
	fmt.Println("int值被修改了，新值为:", i)
}

func moidfy_mapImp(p map[string]int) {
	fmt.Printf("函数里接收到map的内存地址是：%p\n", &p)
	p["zhangsan"] = 20
}
func Modify_Map() {
	persons := make(map[string]int)
	persons["zhangsan"] = 19

	mp := &persons
	fmt.Printf("oringal map addr=%p\n", mp)
	moidfy_mapImp(persons)
	fmt.Println("new map addr=", persons)
}
