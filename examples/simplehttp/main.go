package main

import (
	"log"
	"net/http"
)

func main() {
	handler1 := func(w http.ResponseWriter, r *http.Request) {
		log.Println("handler1が呼び出されました")
	}
	handler2 := func(w http.ResponseWriter, r *http.Request) {
		log.Println("handler2が呼び出されました")
	}

	multiplexer := http.NewServeMux()
	multiplexer.HandleFunc("/", handler1)
	multiplexer.HandleFunc("/abc", handler2)

	if err := http.ListenAndServe(":8080", multiplexer); err != nil {
		// http.ListenAndServeは常にnon-nil errorを返す。
		log.Println(err)
	}
}
