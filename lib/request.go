package lib

import (
	"fmt"
	"net/http"
)

type Request struct {
	http_req *http.Request
	output   chan http.Response
}

func NewRequest(http_req *http.Request) Request {
	return Request{
		http_req: http_req,
		output:   make(chan http.Response),
	}
}

func (r *Request) String() string {
	return fmt.Sprintf("%v", r.http_req.URL)
}
