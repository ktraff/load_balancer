package main

import (
	"fmt"
	"github.com/ktraff/load_balancer/lib"
	"log"
	"net/http"
)

func main() {
	work_channel := make(chan *lib.Request)
	balancer := lib.NewBalancer(10, 5)
	go balancer.Balance(work_channel)

	// Start an HTTP listener that handles all traffic
	h1 := func(_ http.ResponseWriter, http_req *http.Request) {
		fmt.Println(fmt.Sprintf("Received request: %v", http_req))
		req := lib.NewRequest(http_req)
		work_channel <- &req
	}
	http.HandleFunc("/", h1)
	fmt.Println("Serving on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}