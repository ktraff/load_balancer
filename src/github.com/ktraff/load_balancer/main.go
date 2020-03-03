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
    // All workers in the pool will subscribe to this channel to process incoming requests
    work_channel := make(chan *lib.Request)
    // Configure the number of concurrent workers
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
	h1 := func(w http.ResponseWriter, http_req *http.Request) {
		fmt.Println(fmt.Sprintf("Incoming request: %v", http_req))
		req := lib.NewRequest(http_req, w)
		work_channel <- &req
		req.Respond()
	}
	http.HandleFunc("/", h1)
	fmt.Println("Serving on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
