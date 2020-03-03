package lib

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Worker struct {
	id int
	client http.Client    // The client used to proxy requests
	backend url.URL       // The web service to proxy requests to
    requests chan Request // work to do (buffered channel)
    pending  int          // count of pending tasks
    index    int          // index in the heap
}

func NewWorker(id, reqBufferSize int, backend url.URL) Worker {
	return Worker {
		id: id,
		client: http.Client {
			Timeout: 60 * time.Second,
		},
		backend: backend,
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
		go w.forward_request(req.http_req, req.output, done)
    }
}

func (w *Worker) forward_request(req *http.Request, requestor_chan chan http.Response, worker_done_chan chan *Worker) {

	// Recover from unsuccessful requests
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovering from error while forwarding request: %v \n", r)
		}
	}()
	
	url := fmt.Sprintf("%v/%v", w.backend.String(), req.URL.Path)
	fmt.Println(fmt.Sprintf("Forwarding request: %v", url))
	outgoing_req, err := http.NewRequest(req.Method, url, strings.NewReader(""))
	if err != nil {
		fmt.Println(fmt.Sprintf("Error while creating forwarded request (%v): %v", req, err))
	}
	resp, err := w.client.Do(outgoing_req)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error while forwarding request (%v): %v", req, err))
	} else {
		requestor_chan <- *resp
	}

	worker_done_chan <- w
}

func (w Worker) String() string {
	return fmt.Sprintf("Worker #%v => %v", w.id, w.backend.String())
}
