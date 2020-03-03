package main

import (
	"fmt"
	"github.com/ktraff/load_balancer/lib"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	work_channel := make(chan *lib.Request)
	num_workers := 10
	num_requests_per_worker := 5
	if len(os.Args) > 1 {
		num_workers, _ = strconv.Atoi(os.Args[1])
	}
	if len(os.Args) > 2 {
		num_requests_per_worker, _ = strconv.Atoi(os.Args[2])
	}

	balancer := lib.NewBalancer(num_workers, num_requests_per_worker)
	go balancer.Balance(work_channel)

	// Start an HTTP listener that handles all traffic
	h1 := func(_ http.ResponseWriter, http_req *http.Request) {
		fmt.Println(fmt.Sprintf("Incoming request: %v", http_req))
		req := lib.NewRequest(http_req)
		work_channel <- &req
	}
	http.HandleFunc("/", h1)
	fmt.Println("Serving on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
