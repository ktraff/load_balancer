package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	id := os.Getenv("ID")
	port := 8080
	if port_env, exists := os.LookupEnv("PORT"); exists {
		port, _ = strconv.Atoi(port_env)
	}

	handler := func(w http.ResponseWriter, http_req *http.Request) {
		msg := fmt.Sprintf("Hello from %v", id)
		fmt.Println(msg)
		io.WriteString(w, msg)
	}
	http.HandleFunc("/", handler)
	fmt.Println(fmt.Sprintf("Serving on port %v", port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}
