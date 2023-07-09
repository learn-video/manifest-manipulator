package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	exit := make(chan bool)

	handler := http.NewServeMux()
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	go func() {
		http.ListenAndServe(":8080", handler)
	}()

	targetURL := "https://test-streams.mux.dev"
	target, err := url.Parse(targetURL)
	if err != nil {
		log.Fatal(err)
	}

	director := func(req *http.Request) {
		req.URL.Scheme = "https"
		req.URL.Host = target.Host
		req.Host = target.Host
		req.Header.Set("Host", target.Host)
	}
	modifyResponse := func(res *http.Response) error {
		res.Header.Set("Access-Control-Allow-Origin", "*")
		res.Header.Set("Access-Control-Allow-Methods", "*")
		res.Header.Set("Access-Control-Allow-Headers", "*")
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		res.Body = io.NopCloser(bytes.NewReader(body))
		return nil
	}
	proxy := &httputil.ReverseProxy{
		Director:       director,
		ModifyResponse: modifyResponse,
	}

	go func() {
		http.ListenAndServe(":9090", proxy)
	}()

	<-exit
}
