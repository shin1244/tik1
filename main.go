package main

import (
	"net/http"

	"mongo/tiktekto/app.go"

	"github.com/gorilla/pat"
	"github.com/urfave/negroni"
)

func main() {
	mux := pat.New()

	mux.Get("/ws", app.WsHandler)

	go app.HandleTik()

	n := negroni.Classic()
	n.UseHandler(mux)

	http.ListenAndServe(":3000", n)
}
