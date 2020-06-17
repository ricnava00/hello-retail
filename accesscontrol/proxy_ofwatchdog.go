package accesscontrol

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type ProxyServer struct {
	server *http.Server
	port   int
}

func NewProxyServer(port int) ProxyServer {
	mux := http.NewServeMux()
	proxyServer := ProxyServer{
		server: &http.Server{
			Addr:         fmt.Sprintf(":%d", port),
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			Handler:      mux,
		},
		port: port,
	}

	mux.HandleFunc("/", proxyServer.proxy)

	return proxyServer
}

func (s *ProxyServer) Serve() {
	log.Printf("Proxy listening on port: %d\n", s.port)

	go func() {
		err := s.server.ListenAndServe()
		if err != http.ErrServerClosed {
			panic(fmt.Sprintf("proxy error ListenAndServe: %v\n", err))
		}
	}()
}

func (s *ProxyServer) proxy(w http.ResponseWriter, r *http.Request) {
	log.Printf("Starting proxy request to gateway")

	header, err := untrackRequest(r)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = wrapRequestBodyWithHeader(header, r)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	r.Header.Set("Authorization", "service")

	r.URL.Scheme = "http"
	r.URL.Host = "gateway.openfaas.svc.cluster.local:8080"
	r.Host = "gateway.openfaas.svc.cluster.local:8080"

	res, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer res.Body.Close()

	for key, values := range res.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(res.StatusCode)
	io.Copy(w, res.Body)

	log.Printf("Ending proxy request to gateway")
}
