// Package routeutil exposes a route class builder that handles
// different kinds of HTTP requests with method validation.

package routeutil

import (
    "net/http"
    "log"
)

// A route binds a function to one or more HTTP methods.
type Route interface {
    MakeHandler() http.HandlerFunc

}

// A single-method endpoint only accepts one method.
type singleMethodRoute struct {
    method string
    handler func(http.ResponseWriter, *http.Request)
}

func newSingleRoute(method string, fn func(http.ResponseWriter, *http.Request)) Route {
    return &singleMethodRoute {
        method: method,
        handler: fn}
}

// Validates the request method before executing a function.
func (smr *singleMethodRoute) MakeHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != smr.method {
            log.Printf("[ERROR]: HTTP method %s not supported", r.Method)
            w.WriteHeader(http.StatusMethodNotAllowed)
            return

        }
        smr.handler(w, r)
    }
}

// A multi-method route supports multiple methods.
type multiMethodRoute struct {
    methods []string
    handler func(http.ResponseWriter, *http.Request)
}

func newMultiRoute(methods []string, handler func(http.ResponseWriter, *http.Request)) Route {
    return &multiMethodRoute{
        methods: methods,
        handler: handler}
}

// Validate request method against a list of acceptable methods
func (mmr *multiMethodRoute) MakeHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        validMethod := inList(r.Method, mmr.methods)
        if ! validMethod  {
            log.Printf("[ERROR]: HTTP method %s not supported", r.Method)
            w.WriteHeader(http.StatusMethodNotAllowed)
            return
        }

        mmr.handler(w, r)
    }
}

// Checks to see if a value is in a list of choices
func inList(val string, choices []string) bool {
    for i, j := 0, len(choices); i < j; i++ {
        if choices[i] == val {
            return true
        }
    }
    return false
}

// Routes that only support HTTP GET
func GETRoute(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
    return newSingleRoute(http.MethodGet, fn).MakeHandler()
}

// POST routes must support POST and OPTIONS (for CORS)
func POSTRoute(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
    return newMultiRoute([]string{http.MethodPost, http.MethodOptions}, fn).MakeHandler()
}
