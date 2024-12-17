package main

import (
	"crypto/tls"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

func handleTunneling(w http.ResponseWriter, r *http.Request) {
	destConn, err := net.DialTimeout("tcp", r.Host, 10*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}
	go transfer(destConn, clientConn)
	go transfer(clientConn, destConn)
}

func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}

func handleHTTP(w http.ResponseWriter, req *http.Request) {

	//req.URL.Scheme = "http"
	s := strings.Split(req.Host, ":")
	addr := net.ParseIP(s[0])
	if addr != nil {
		req.URL.Scheme = "http"
	} else {
		req.URL.Scheme = "https"
	}

	req.URL.Host = req.Host

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	copyHeader(w.Header(), resp.Header)
	w.Header().Add("x-test-header", "12345")

	w.WriteHeader(resp.StatusCode)

	finishRead := false
	for !finishRead {
		b := make([]byte, 1024)
		n, err := resp.Body.Read(b)

		if err == io.EOF {
			finishRead = true
		}

		if n > 0 {
			w.Write(b[:n])
		}
	}

	//bytes := []byte("Hey, I'm taking over this body!")
	//w.Header().Set("Content-Length", strconv.Itoa(len(bytes)))

	// b, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println(string(b))
	// }

	//io.Copy(w, resp.Body)
	//w.Write(bytes)
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func main() {
	server := &http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodConnect {
				handleTunneling(w, r)
			} else {
				handleHTTP(w, r)
			}
		}),
		// Disable HTTP/2.
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	log.Fatal(server.ListenAndServe())
}
