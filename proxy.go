package main

import (
	"fmt"
	"io"
	"net/http"
)

func copyHeaders(dst http.Header, src http.Header) {
	for key, values := range src {
		for _, value := range values {
			dst.Add(key, value)
		}
	}
}

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

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error closing the response body", err)
		}
	}(resp.Body)

	w.WriteHeader(resp.StatusCode)
	copyHeaders(w.Header(), resp.Header)

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		http.Error(w, "Error copying response body", http.StatusInternalServerError)
		return
	}
}
