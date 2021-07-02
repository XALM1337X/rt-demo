package main

import (
    "fmt"
    "html"
    "log"
    "net/http"
)

func main() {

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	fmt.Fprintf(w, "Request sent: %#v",r)
    })

    http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request){
        fmt.Fprintf(w, "Hi")
    })

    log.Fatal(http.ListenAndServe(":6543", nil))

}



/*
static := http.FileServer(http.Dir("./static"))
http.Handle("/",static)
http.ListenAndServe(":8081",nil)
*/
