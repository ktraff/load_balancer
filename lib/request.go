package lib

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type Request struct {
	http_req        *http.Request        // The incoming http request
	output          chan http.Response   // The channel for receiving the response from the worker pool
	response_writer http.ResponseWriter  // The response writer used to communicate back to the client
}

func NewRequest(http_req *http.Request, r http.ResponseWriter) Request {
	return Request{
		http_req:        http_req,
		output:          make(chan http.Response),
		response_writer: r,
	}
}

func (r Request) Respond() {
	// Wait to receive a response from the worker pool
	resp := <-r.output

	// Relay the response back to the client
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
