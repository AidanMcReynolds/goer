package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"syscall/js"
)

func main() {

	fmt.Println("intializing")

	r, err := http.Get("http://localhost/code.go")

	if err != nil {

		js.Global().Get("document").Call("getElementById", "source").Set("innerHTML", err)

	} else {

		b, err := ioutil.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			js.Global().Get("document").Call("getElementById", "source").Set("innerHTML", err)
		} else {
			js.Global().Get("document").Call("getElementById", "source").Set("innerHTML", string(b))
			fmt.Println(string(b))
		}
	}

}
