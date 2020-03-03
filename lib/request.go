package lib

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type Request struct {
	http_req *http.Request
	output   chan http.Response
	response_writer http.ResponseWriter
}

func NewRequest(http_req *http.Request, r http.ResponseWriter) Request {
	return Request{
		http_req: http_req,
		output:   make(chan http.Response),
		response_writer: r,
	}
}

func (r Request) Respond() {
	resp := <-r.output
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
    if err != nil {
        log.Fatal(err)
    }
	bodyString := string(bodyBytes)
	fmt.Println(fmt.Sprintf("Responding to request with: %v", bodyString))
	io.WriteString(r.response_writer, bodyString)
}

func (r *Request) String() string {
	return fmt.Sprintf("%v", r.http_req.URL)
}
