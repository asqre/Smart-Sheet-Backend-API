package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", HelloWorld)         // set router
	err := http.ListenAndServe(":8080", nil) //set listen port
	if err != nil {
		log.Println("ListenAndServe at :8080", err)
	}
}

// TODO :do in JSON APi's
func HelloWorld(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm() // it parses(receive) raw query of all requests from the URL and updates r.Form
	if err != nil {
		return
	}
	fmt.Println(r.Form)
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["URL_Long"])
	for k, v := range r.Form { // The range keyword is used in a for loop to iterate over the elements of an array, slice, map, or string.
		fmt.Println("key", k)
		fmt.Println("value", v)
	}
	_, err = fmt.Fprintf(w, "Welcome to the go programming")
	if err != nil {
		return
	}
}
