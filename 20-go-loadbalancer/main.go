package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type Server interface {
    // Returns the address with which to access the server
    Address() string
    // Returns true if the server is able to serve requests
    IsAlive() bool
    // Serve uses this server to process the request
    Serve(http.ResponseWriter, *http.Request)
}

type simpleServer struct {
    addr string
    proxy *httputil.ReverseProxy
}

func (s *simpleServer) Address() string {
    return s.addr
}

func (s *simpleServer) IsAlive() bool {
    return true
}

func (s *simpleServer) Serve(r http.ResponseWriter, w *http.Request) {
    s.proxy.ServeHTTP(r, w)
}

func newSimpleServer(addr string) *simpleServer {
    serverUrl, err := url.Parse(addr)
    handleErr(err)

    return &simpleServer{
        addr: addr,
        proxy: httputil.NewSingleHostReverseProxy(serverUrl),
    }
}

type LoadBalancer struct {
    port string
    roundRobinCount int
    servers []Server
}

func NewLoadBalancer(port string, servers []Server) *LoadBalancer {
    return &LoadBalancer{
        port: port,
        roundRobinCount: 0,
        servers: servers,
    }
}


// Returns the address of the next available server to send a request to, 
// using a simple Round-Robin algorithm
func (lb *LoadBalancer) getNextAvailableServer() Server {
    server := lb.servers[lb.roundRobinCount % len(lb.servers)]

    // If the current server is busy, then upgrade and check for next one
    for !server.IsAlive(){
        fmt.Println("Server did not respond")
        lb.roundRobinCount++
        server = lb.servers[lb.roundRobinCount % len(lb.servers)]
    }
    lb.roundRobinCount++
    return server

    // Ideally we should be setting up some mechanism where after a server 
    // has completed the task, it if offloaded from the roundRobintCount
}

func (lb *LoadBalancer) serveProxy(w http.ResponseWriter, r *http.Request){
    targetServer := lb.getNextAvailableServer()
    fmt.Printf("Forwarding request to address:%q\n", targetServer.Address())
    // Delete the X-Forwarding-For header to prevent IP spoofing
    targetServer.Serve(w, r)
}

func main(){
    servers := []Server{
        newSimpleServer("https://www.facebook.com"),
        newSimpleServer("https://www.bing.com"),
        newSimpleServer("https://www.duckduckgo.com"),
    }
    lb := NewLoadBalancer("8000", servers)

    handleRedirect := func(w http.ResponseWriter, r *http.Request){
        lb.serveProxy(w, r)
    }

    // register a proxy handler to handle all requests
    http.HandleFunc("/", handleRedirect)

    fmt.Printf("Serving requests at 'localhost:%s'\n", lb.port)
    http.ListenAndServe(":" + lb.port, nil)
}

// Simple error handler
func handleErr(err error) {
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }
}
