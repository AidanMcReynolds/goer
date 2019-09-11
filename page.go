package main

import (
	"net/http"
)

func postPage(w http.ResponseWriter, req *http.Request) {
	m := req.URL.Query()
	token := m.Get("page")
	c := &http.Cookie{Name: "temp", Value: token}
	http.SetCookie(w, c)
	http.ServeFile(w, req, "static/out.html")
}
