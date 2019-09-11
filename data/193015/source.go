package main

import "syscall/js"

func main() {
	js.Global().Get("document").Call("getElementById", "info").Set("innerHTML", "hello")
}
