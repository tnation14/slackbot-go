package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "slackbot-go/util"
)

func echoBody(w http.ResponseWriter, r *http.Request){
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        log.Printf("[ERROR] %s", err)
        return
    }

    w.WriteHeader(200)
    w.Header().Set("Content-Type", "application/json")
    w.Write(body)
}

func sayHello(w http.ResponseWriter, r *http.Request){
    w.WriteHeader(200)
    w.Header().Set("Content-Type", "text/plain")
    fmt.Fprintf(w, "This is the webserver root!")
}


func main() {


    http.HandleFunc("/", routeutil.GETRoute(sayHello))
    http.HandleFunc("/echo", routeutil.POSTRoute(echoBody))
    log.Printf("Starting up")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
