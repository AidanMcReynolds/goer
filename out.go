package main

import (
	"archive/zip"
	"bytes"
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
	token := getToken(w, req)
	err := saveSource(req, token)
	if err != nil {
		io.WriteString(w, err.Error())
	} else {
		err = compile(token)
		if err != nil {
			io.WriteString(w, err.Error())
		} else {

			if req.FormValue("button") == "Run >>>" {
				runPage(w, req)
			} else if req.FormValue("button") == "download" {
				wasmPage(w, req, token)
			}
		}
	}
}

func runPage(w http.ResponseWriter, req *http.Request) {

	http.ServeFile(w, req, "static/out.html")

}
func wasmPage(w http.ResponseWriter, req *http.Request, t string) {
	buf := new(bytes.Buffer)
	zipw := zip.NewWriter(buf)
	zipAdd("out.wasm", "data/"+t+"/out.wasm", zipw)
	zipAdd("out.html", "static/out.html", zipw)
	zipAdd("wasm_exec.js", "static/wasm_exec.js", zipw)
	zipw.Close()
	http.ServeContent(w, req, "your.wasm", time.Now(), bytes.NewReader(buf.Bytes()))

}

func zipAdd(name string, path string, w *zip.Writer) error {
	f, err := w.Create(name)
	if err != nil {
		return err
	}
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, file)
	return err
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
	cmd := exec.Command("/tmp/go/bin/go", "build", "-o", "data/"+token+"/out.wasm", "data/"+token+"/source.go") //
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
