// +build ignore

package main

import "unsafe"

//import (
//	"fmt"
//	"unsafe"
//)
//
//type Num struct{
//	i string
//	j int64
//}
//
//func main(){
//	n := Num{i: "EDDYCJY", j: 1}
//	nPointer := unsafe.Pointer(&n)
//	fmt.Printf("n.i: %s, n.j: %d\n", n.i, n.j)
//	niPointer := (*string)(unsafe.Pointer(nPointer))
//	*niPointer = "煎鱼"
//
//	njPointer := (*int64)(unsafe.Pointer(uintptr(nPointer) + unsafe.Offsetof(n.j)))
//	*njPointer = 2
//
//	fmt.Printf("n.i: %s, n.j: %d", n.i, n.j)
//}

func main() {
	n := 5
	np := unsafe.Pointer(&n)
	//nf := (*float64)(unsafe.Pointer(&n))
	println(np)
	//println(nf)
	println(uintptr(unsafe.Pointer(&n)))
}