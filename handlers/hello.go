package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type Hello struct {
	logger *log.Logger
}

func NewHello(logger *log.Logger) *Hello {
	return &Hello{logger}
}

func (hello *Hello) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	hello.logger.Println("Hello World")
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "Opps! something went wrong", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(rw, "Hello %s\n", data)
}
