package main

import "fmt"


type Base interface {
	Input() int
}

type Dog struct {
}

func (p Dog) Input() int {
	fmt.Println("call input dog")
	return 100
}

func GetDog() Base {
	return Dog{}
}

func tst_Dog() {
	// handler := GetDog()
	// handler := func() Base { return Dog() }
	handler := func() Base { return Dog{} }

	fmt.Printf("type: %T, %#v	\n", handler(), handler().Input())

}

type USB interface {
	Name() string
	Connecter //嵌入Connecter，从而USB就拥有Connecter的方法Connect()
}

type Connecter interface {
	Connect() int
}

type Phone struct {
	phone_name string
}

func (cobj Phone) Name() string {
	return cobj.phone_name
}

func (cobj Phone) Connect() int {
	fmt.Println(cobj.phone_name, " connect successed!")
	return 0
}

func Disconnect(oneI interface{}) { // 控接口参数，可以传入任何对象
	//当有空接口时，一般选用switch语句，使用type switch则可针对空接口进行比较全面的类型判断
	switch v := oneI.(type) {
	case Phone:
		fmt.Println("aaaaaaaaaaaaaaaa")
		fmt.Println(v.Name(), " is Disconnected")
	default:
		fmt.Println("未知设备!!!")
	}
}

func tst_phone() {

	var pUsbPtr USB
	pUsbPtr = Phone{"test obj"}
	pUsbPtr.Connect()

	Disconnect(pUsbPtr)
}

func tst_fun_entry() {

	tst_I_ex()
	return

	tst_phone()
	return

	tst_Dog()

}
