package main

import (
	"net/http"
)

func main() {
	InitAuth()
	chanquit := make(chan int)
	server := Config.Server()
	mux := http.NewServeMux()
	mux.Handle(Auth.Path+"/", http.StripPrefix(Auth.Path, WrapRecover(
		Session.Wrap(
			http.HandlerFunc(
				Auth.Serve(ActionSuccess),
			),
		),
	)),
	)
	mux.Handle("/", WrapRecover(
		http.HandlerFunc(
			ActionIndex,
		),
	),
	)
	server.Handler = mux

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
	<-chanquit
}
