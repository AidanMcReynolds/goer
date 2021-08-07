package main

import (
	"strconv"
	"syscall/js"
)

func main() {
	html := `Elements of the fibonacci sequence from the 0th to the nth element.</br><input type="text" id="n"/><span id="err"></span></br>
    <button onClick="fib();">run</button></br>
    <div id="result"></div>`
	js.Global().Get("document").Call("getElementById", "info").Set("innerHTML", html)

	c := make(chan struct{}, 0)
	js.Global().Set("fib", js.FuncOf(fib))
	<-c
}

func fib(this js.Value, args []js.Value) interface{} {

	val := js.Global().Get("document").Call("getElementById", "n").Get("value").String()
	n, err := strconv.Atoi(val)
	if n < 0 || err != nil {
		js.Global().Get("document").Call("getElementById", "err").Set("innerHTML", " invalid input")
	} else {
		js.Global().Get("document").Call("getElementById", "err").Set("innerHTML", "")
		js.Global().Get("document").Call("getElementById", "result").Set("innerHTML", fibn(n))
	}
	return nil

}
func fibn(n int) string {
	x := []int{0}

	if n > 0 {
		x = append(x, 1)
		if n > 1 {
			for i := 2; i <= n; i++ {
				x = append(x, x[i-1]+x[i-2])
			}
		}
	}

	//convert to string
	y := ""
	for _, v := range x {
		y = y + " " + strconv.Itoa(v)
	}
	return y
}
