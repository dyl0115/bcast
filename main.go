package main

import (
	"log"
	"net/http"
)

func main() {
	go InitInjectQueueWorker()
	http.HandleFunc("/inject", handleInject)
	http.HandleFunc("/stream", handleStream)
	log.Fatal(http.ListenAndServe(":8336", nil))
}
