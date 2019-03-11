// Package routeutil exposes a route class builder that handles
// different kinds of HTTP requests with method validation.

package routeutil

import (
    "net/http"
    "log"
)

// A route wraps a HTTP handler in a validator function.
type Route interface {
    MakeHandler() http.HandlerFunc

}

type routeImpl struct {
    validator func(*http.Request) (bool, int)
    handler func(http.ResponseWriter, *http.Request)
}

func newRoute(fn func(http.ResponseWriter, *http.Request), validator func(*http.Request) (bool, int)) Route {
    return &routeImpl {
        handler: fn,
        validator: validator}
}

// Validates the request method before executing a function.
func (route *routeImpl) MakeHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        valid, errcode := route.validator(r)
        if valid  {
            route.handler(w, r)
        } else {
            w.WriteHeader(errcode)
            return
        }
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
    validator := func(r *http.Request) (bool, int){
                    if r.Method == http.MethodGet {
                        return true, http.StatusOK
                    }
                    log.Printf("[ERROR] Method %s not allowed", r.Method)
                    return false, http.StatusMethodNotAllowed
                 }
    return newRoute(
            fn,
            validator).MakeHandler()
}


func POSTRoute(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
    validator := func(r *http.Request) (bool, int) {
                    if !inList(r.Method, []string{http.MethodPost, http.MethodOptions}) {
                        log.Printf("[ERROR]: Method %s not in %s.", r.Method, []string{http.MethodPost, http.MethodOptions})
                        return false, http.StatusMethodNotAllowed
                    }
                        return true, http.StatusOK
                }
    return newRoute(
            fn,
            validator).MakeHandler()
}

// Wraps a route in a function that checks username and password
func AuthenticatedRoute(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
    // TODO Make this validator real
    validator := func(r *http.Request) (bool, int) {
                    log.Printf("[INFO]: Username or password incorrect.")
                    return false, http.StatusForbidden
                }
    return newRoute(
        fn,
        validator).MakeHandler()
}

func AuthenticatedPOSTRoute (fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
    return AuthenticatedRoute(GETRoute(fn))
}
