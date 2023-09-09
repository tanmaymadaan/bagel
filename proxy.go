package main

import (
	"io"
	"net/http"
)

type Proxy struct {
	targetUrl string
}

func NewProxy(targetUrl string) *Proxy {
	return &Proxy{targetUrl: targetUrl}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	target := p.targetUrl + r.URL.String()

	resp, err := http.Get(target)
	if err != nil {
		http.Error(w, "Error making request to backend server", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	copyHeaders(w.Header(), resp.Header)

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		http.Error(w, "Error copying response body", http.StatusInternalServerError)
		return
	}
}