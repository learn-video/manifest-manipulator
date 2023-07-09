package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/learn-video/manifest-manipulator/filter"
)

func modifyManifest(res *http.Response) (io.Reader, error) {
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if strings.HasSuffix(res.Request.URL.Path, "/master.m3u8") {
		p, err := filter.NewMasterPlaylist(*bytes.NewBuffer(body))
		if err != nil {
			return nil, err
		}
		p.FilterBandwidth(filter.BandwidthFilter{Max: 3000000})

		body = []byte(p.Playlist.String())
	}
	return bytes.NewReader(body), nil
}

func main() {
	exit := make(chan bool)

	handler := http.NewServeMux()
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	go func() {
		http.ListenAndServe(":8080", handler)
	}()

	targetURL := "https://cph-p2p-msl.akamaized.net"
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

		body, err := modifyManifest(res)
		if err != nil {
			return err
		}

		res.Body = io.NopCloser(body)
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
