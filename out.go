package main

import (
	"errors"
	"go/format"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func outputPage(w http.ResponseWriter, req *http.Request) {

	if req.FormValue("button") == "Run >>>" {
		runPage(w, req)
	} else if req.FormValue("button") == "download WASM" {
		wasmPage(w, req)
	} else if req.FormValue("button") == "download source" {
		dlPage(w, req)
	} else if req.FormValue("button") == "link project" {
		projectPage(w, req)
	} else if req.FormValue("button") == "link output" {
		userPage(w, req)
	}
}

func projectPage(w http.ResponseWriter, req *http.Request) {
	token := getToken(req)
	err := saveSource(req, token)
	if err != nil {
		io.WriteString(w, err.Error())
	} else {
		http.Redirect(w, req, "http://localhost/?user="+token, 307)
	}
}
func userPage(w http.ResponseWriter, req *http.Request) {
	token := getToken(req)
	err := saveSource(req, token)
	if err != nil {
		io.WriteString(w, err.Error())
	} else {
		err = compile(token)
		if err != nil {
			io.WriteString(w, err.Error())
		} else {
			http.Redirect(w, req, "http://localhost/p/?page="+token, 307)
		}
	}
}

func runPage(w http.ResponseWriter, req *http.Request) {
	token := getToken(req)
	err := saveSource(req, token)
	if err != nil {
		io.WriteString(w, err.Error())
	} else {
		err = compile(token)
		if err != nil {
			io.WriteString(w, err.Error())
		} else {
			http.ServeFile(w, req, "static/out.html")
		}
	}

}
func wasmPage(w http.ResponseWriter, req *http.Request) {
	token := getToken(req)
	err := saveSource(req, token)
	if err != nil {
		io.WriteString(w, err.Error())

	} else {
		err = compile(token)
		if err != nil {
			io.WriteString(w, err.Error())
		} else {
			f, _ := os.Open("data/" + token + "/out.wasm")
			http.ServeContent(w, req, "your.wasm", time.Now(), f)
			f.Close()
		}
	}
}
func dlPage(w http.ResponseWriter, req *http.Request) {
	token := getToken(req)
	err := saveSource(req, token)
	if err != nil {
		io.WriteString(w, err.Error())
	} else {
		f, _ := os.Open("data/" + token + "/source.go")
		http.ServeContent(w, req, "your.go", time.Now(), f)
		f.Close()
	}
}
func saveSource(req *http.Request, token string) error {
	source := []byte(req.FormValue("source"))
	source, err := format.Source(source)
	if err == nil {
		err = ioutil.WriteFile("data/"+token+"/source.go", source, 0644)

	}
	return err
}

func compile(token string) error {
	cmd := exec.Command("/tmp/go/bin/go", "build", "-o", "data/"+token+"/out.wasm", "data/"+token+"/source.go")
	cmd.Env = append(os.Environ(), "GOARCH=wasm", "GOOS=js")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(err.Error() + "   " + string(out))
	}
	if len(out) > 0 {
		return errors.New(string(out))
	}
	return nil
}
