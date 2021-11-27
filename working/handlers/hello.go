package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Hello struct {
	l *log.Logger
}

func NewHello ( l *log.Logger) *Hello {
	return &Hello{l}
}
func (h *Hello) ServeHTTP(rw http.ResponseWriter, r *http.Request){

		h.l.Println("Hello World")
		//reading from request body
		d, err := ioutil.ReadAll(r.Body)

		//error catching
		if err != nil {
			http.Error(rw, "Oops", http.StatusBadRequest)
			return
			//Alternative.
			// rw.WriteHeader(http.StatusBadRequest)
			// rw.Write([]byte("Oops!"))
		}
		
		fmt.Fprintf(rw, "Hello %s", d)
}