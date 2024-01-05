package main

import (
	"log"
	"net/http"
)

func main() {
	// runWithoutPackageVariable()
	// runWithoutMultiplexer()
	runWithPackageVariables()
}

// ServerとServeMuxを関数内で作るタイプの実装
func runWithoutPackageVariable() {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/", handler1)
	serveMux.HandleFunc("/abc", handler2)
	serveMux.HandleFunc("/favicon.ico", handlerForFavicon)

	server := http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}

	if err := server.ListenAndServe(); err != nil {
		// (*http.Server).ListenAndServeは常にnon-nil errorを返す。
		log.Println(err)
	}
}

func runWithoutMultiplexer() {
	server := http.Server{
		Addr:    ":8080",
		Handler: http.HandlerFunc(handler1),
	}

	if err := server.ListenAndServe(); err != nil {
		// (*http.Server).ListenAndServeは常にnon-nil errorを返す。
		log.Println(err)
	}
}

func runWithPackageVariables() {
	// パスが/にマッチするリクエストはhandler1にルーティングする
	http.HandleFunc("/", handler1)
	// パスが/abcにマッチするリクエストはhandler2にルーティングする
	http.HandleFunc("/abc", handler2)
	// パスが/favicon.icoにマッチするリクエストはfaviconhandlerにルーティングする
	http.HandleFunc("/favicon.ico", handlerForFavicon)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		// http.ListenAndServeは常にnon-nil errorを返す。
		log.Println(err)
	}
}

func handler1(w http.ResponseWriter, r *http.Request) {
	log.Println("handler1が呼び出されました")
	w.Write([]byte("I'm handler 1"))
}

func handler2(w http.ResponseWriter, r *http.Request) {
	log.Println("handler2が呼び出されました")
	w.Write([]byte("I'm handler 2"))
}

// Webブラウザーは/favicon.icoにリクエストを送るのでそのためのハンドラーも別に作っておく
func handlerForFavicon(w http.ResponseWriter, r *http.Request) {
	log.Println("favicon handlerが呼び出されました")
	// w.Write([]byte("I'm handler 2")) TODO: gopherくん画像を返す
}
