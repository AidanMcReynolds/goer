package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"
)

func main() {
	println("starting server . . . . . . ")
	port, err := getPort()
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", mainPage)
	mux.HandleFunc("/out/", outputPage)
	mux.HandleFunc("/p/", postPage)
	mux.HandleFunc("/out/wasm_exec.js", wasmexecjs)
	mux.HandleFunc("/p/wasm_exec.js", wasmexecjs)
	mux.HandleFunc("/out/out.wasm", outwasm)
	mux.HandleFunc("/p/out.wasm", poutwasm)
	err = http.ListenAndServe(port, mux)
	fmt.Println(err)

}
func getPort() (string, error) {
	// the PORT is supplied by Heroku
	port := os.Getenv("PORT")
	if port == "" {
		return "", fmt.Errorf("$PORT not set")
	}
	return ":" + port, nil
}
func wasmexecjs(w http.ResponseWriter, req *http.Request) {
	f, _ := os.Open("static/wasm_exec.js")
	http.ServeContent(w, req, "static/wasm_exec.js", time.Now(), f)
	f.Close()
}

func outwasm(w http.ResponseWriter, req *http.Request) {
	token := getToken(req)
	f, _ := os.Open("data/" + token + "/out.wasm")
	http.ServeContent(w, req, "out.wasm", time.Now(), f)
	f.Close()
}
func poutwasm(w http.ResponseWriter, req *http.Request) {
	c, err := req.Cookie("temp")
	if err != nil {
		fmt.Println(err)
	}
	token := c.Value

	f, err := os.Open("data/" + token + "/out.wasm")
	if err != nil {
		fmt.Println(err)
	}
	http.ServeContent(w, req, "out.wasm", time.Now(), f)
	f.Close()
}

func mainPage(w http.ResponseWriter, req *http.Request) {
	var t string
	m := req.URL.Query()
	if m.Get("user") != "" {
		t = m.Get("user")

		c := &http.Cookie{Name: "auth", Value: t}
		http.SetCookie(w, c)
	} else {
		t = getToken(req)
		if t == "" {
			c := &http.Cookie{Name: "auth", Value: newToken()}
			http.SetCookie(w, c)
		}
	}

	tmpl, _ := template.ParseFiles("static/index.html")
	f, err := ioutil.ReadFile("data/" + t + "/source.go")
	if err != nil {
		f, _ = ioutil.ReadFile("static/default.go")
	}
	s := sourcefile{Source: string(f)}
	tmpl.Execute(w, s)
}

func getToken(req *http.Request) string {
	c, err := req.Cookie("auth")

	if err == nil {
		return c.Value
	}
	return ""
}

func newToken() string {
	for {
		i := strconv.Itoa(rand.Intn(1000000))
		_, err := os.Stat("data/" + i)
		if os.IsNotExist(err) {
			// doesn't exist
			os.Mkdir("data/"+i, 0777) //0777 is permenant public file
			return i
		}
	}

}

type sourcefile struct {
	Source string
}
