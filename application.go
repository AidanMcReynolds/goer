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

	mux.HandleFunc("/out/wasm_exec.js", wasmexecjs)

	mux.HandleFunc("/out/out.wasm", outwasm)

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
	token := getToken(w, req)
	f, _ := os.Open("data/" + token + "/out.wasm")
	http.ServeContent(w, req, "out.wasm", time.Now(), f)
	f.Close()
}

func mainPage(w http.ResponseWriter, req *http.Request) {

	t := getToken(w, req)

	tmpl, _ := template.ParseFiles("static/index.html")
	f, err := ioutil.ReadFile("data/" + t + "/source.go")
	if err != nil {
		f, _ = ioutil.ReadFile("static/default.go")
	}
	s := sourcefile{Source: string(f)}
	tmpl.Execute(w, s)
}

func getToken(w http.ResponseWriter, req *http.Request) string {
	c, err := req.Cookie("auth")

	if err != nil {
		c = &http.Cookie{Name: "auth", Value: newToken()}
		http.SetCookie(w, c)
	}
	//check if foler exists

	_, err = os.Stat("data/" + c.Value)
	if os.IsNotExist(err) {
		c = &http.Cookie{Name: "auth", Value: newToken()}
		http.SetCookie(w, c)
	}
	return c.Value

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
