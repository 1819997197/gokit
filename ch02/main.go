package main

import (
	"fmt"
	"gokit/server"
	"net/http"
)

func main() {
	var service server.OrderServer
	service = &server.OrderService{}
	server.InitLog()
	service = server.LoggingMiddleware(server.Log)(service)
	handler := server.MakeHTTPHandler(service)
	fmt.Println("listen:50050")
	err := http.ListenAndServe(":50050", handler)
	if err != nil {
		fmt.Println("http.ListenAndServe err: ", err)
	}
}
