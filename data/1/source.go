package main

import (
	"fmt"
	"syscall/js"
	"time"
)

func main() {
	page := "hello <br/><div id='output'><div>"
	fmt.Println("hello, world")
	js.Global().Get("document").Call("getElementById", "info").Set("innerHTML", page)
	time.Sleep(time.Second)
	js.Global().Get("document").Call("getElementById", "output").Set("innerHTML", "hello")
}
