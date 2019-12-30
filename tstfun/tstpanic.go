package main

import (
	"fmt"
)

type Handler interface {
	Filter(err error, r interface{}) error
}

type Logger interface {
	Ef(format string, a ...interface{})
}

// Handle panic by hdr, which filter the error.
// Finally log err with logger.
func HandlePanic(hdr Handler, logger Logger) error {
	return handlePanic(recover(), hdr, logger)
}

type hdrFunc func(err error, r interface{}) error

func (v hdrFunc) Filter(err error, r interface{}) error {
	return v(err, r)
}

type loggerFunc func(format string, a ...interface{})

func (v loggerFunc) Ef(format string, a ...interface{}) {
	v(format, a...)
}

// Handle panic by hdr, which filter the error.
// Finally log err with logger.
func HandlePanicFunc(hdr func(err error, r interface{}) error, logger func(format string, a ...interface{})) error {
	var f Handler
	if hdr != nil {
		f = hdrFunc(hdr)
	}

	var l Logger
	if logger != nil {
		l = loggerFunc(logger)
	}

	return handlePanic(recover(), f, l)
}

func handlePanic(r interface{}, hdr Handler, logger Logger) error {
	if r != nil {
		err, ok := r.(error)
		if !ok {
			err = fmt.Errorf("r is %v", r)
		}

		if hdr != nil {
			err = hdr.Filter(err, r)
		}

		if err != nil && logger != nil {
			logger.Ef("panic err %+v", err)
		}

		return err
	}

	return nil
}

// IT干货栈 定义一个结构体[大小写敏感]
type Car struct {
	Name   string  // 名称
	Color  string  // 颜色
	Length float32 // 长度
}

// 定义一个小车 结构
type SmallCar struct {
	Car            // 车
	Height float32 // 高度

}

func (car *Car) run() {
	fmt.Println(car.Name, "正在迅速行驶。。。。")
}

func (car *Car) fly() {
	fmt.Println(car.Name, "正在飞行。。。。")
}
func (car Car) changeName() {
	car.Name = "保时捷"
}
func (car *Car) realChangeName() {
	car.Name = "宝马"
}
func TstCar() {
	var car Car
	car.Name = "小栈"
	car.Color = "red"
	car.Length = 2.0
	fmt.Println(car)
	fmt.Printf("Name: %p\n", &car.Name)
	fmt.Printf("Color: %p\n", &car.Color)
	fmt.Printf("Length: %p\n", &car.Length)
}

func TstPanic() {
	TstCar()
	return
	func() {
		defer HandlePanicFunc(nil, func(format string, a ...interface{}) {
			fmt.Println(fmt.Sprintf(format, a...))
		})
		panic("ok 111")
	}()
	// logger := func(format string, a ...interface{}) {
	// 	fmt.Println(fmt.Sprintf(format, a...))
	// }
	// func() {
	// 	defer HandlePanicFunc(nil, logger)
	// 	panic("ok 222")
	// }()
}
