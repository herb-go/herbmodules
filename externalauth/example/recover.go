package main

import "net/http"

func WrapRecover(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			r := recover()
			if r != nil {
				err := r.(error)
				w.WriteHeader(500)
				w.Write([]byte(err.Error()))
			}
		}()
		handler.ServeHTTP(w, r)
	})
}
