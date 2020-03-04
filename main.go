package main

import (
	"log"
	"net/http"

	"github.com/oarfah/codis-scheduler/handle"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte("ok")) })
	http.HandleFunc("/predicate", handle.PredicateHandler)
	http.HandleFunc("/prioritize", handle.PrioritizeHandler)
	log.Panic(http.ListenAndServe(":8880", nil))
}
