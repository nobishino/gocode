package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

func handler(w http.ResponseWriter, r *http.Request) {
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	fmt.Println(string(dump))
	fmt.Fprintf(w, "<html><body>hello</body></html>\n")
}

func main() {
	var httpServer http.Server
	http.HandleFunc("/", handler)

	// well known port 以外が2つに分かれるかも
	log.Println("start http listening :18888")
	httpServer.Addr = ":18888"
	log.Println(httpServer.ListenAndServe())
}

//  serverを使うメリットとして、signal.Notifyをつかったgraceful shutdownがある
// nc -- arbitrary TCP and UDP connections and listens
// https://pkg.go.dev/net/textproto
