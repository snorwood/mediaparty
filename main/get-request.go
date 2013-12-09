package main

import (
	"fmt"
	"net/http"
)

func mp3Handler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "C:\\Users\\Public\\Music\\Sample Music\\Kalimba.mp3")
	fmt.Println("Served")
}

func musicHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `{\"artist\": \"MOOOOO\"}`)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "request.html")
}

func main() {
	http.HandleFunc("/mp3/", mp3Handler)
	http.HandleFunc("/music/", musicHandler)
	http.HandleFunc("/", rootHandler)
	http.ListenAndServe(":8000", nil)
}
