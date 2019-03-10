package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
)

func rootHandler(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "This is the webserver root!")
}

func makePostHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost && r.Method != http.MethodOptions {
            log.Printf("[ERROR]: HTTP method %s not supported", r.Method)
            w.WriteHeader(http.StatusMethodNotAllowed)
            return

        }
        fn(w, r)
    }
}

func makeGetHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet {
            log.Printf("[ERROR]: HTTP method %s not supported", r.Method)
            w.WriteHeader(http.StatusMethodNotAllowed)
            return

        }
        fn(w, r)
    }
}

func echoHandler(w http.ResponseWriter, r *http.Request){
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        log.Printf("[ERROR] %s", err)
        return
    }

    w.WriteHeader(200)
    w.Header().Set("Content-Type", "application/json")
    w.Write(body)
}


func main() {

    http.HandleFunc("/", makeGetHandler(rootHandler))
    http.HandleFunc("/echo", makePostHandler(echoHandler))
    log.Fatal(http.ListenAndServe(":8080", nil))
}
