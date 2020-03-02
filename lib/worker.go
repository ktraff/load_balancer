package lib

import (
	"fmt"
	"net/http"
	"time"
)

type Worker struct {
	client http.Client
    requests chan Request // work to do (buffered channel)
    pending  int          // count of pending tasks
    index    int         // index in the heap
}

func NewWorker(reqBufferSize int) Worker {
	return Worker {
		client: http.Client {
			Timeout: 60 * time.Second,
		},
		requests: make(chan Request, reqBufferSize),
		pending: 0,
		index: 0,
	}
}

func (w *Worker) work(done chan *Worker) {
	fmt.Println("Beginning work")
    for {
		req := <-w.requests
		fmt.Println(fmt.Sprintf("Received a request: %v", req.http_req.URL.String()))
		resp := w.forward_request(req.http_req)
		req.output <- *resp
		done <- w
    }
}

func (w *Worker) forward_request(req *http.Request) *http.Response {
	fmt.Println(fmt.Sprintf("Forwarding request: %v %v", req.Method, req.URL.String()))
	outgoing_req, err := w.client.NewRequest("GET", "http://google.com")
	if err != nil {
		fmt.Println(fmt.Sprintf("Error while creating forwarded request (%v): %v", req, err))
	}
	resp, err := w.client.Do(outgoing_req)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error while formatting request (%v): %v", req, err))
	}
	return resp
}