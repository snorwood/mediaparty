package main

import (
	"github.com/snorwood/mediaparty"
	"net/http"
)

func main() {
	// Define handlers
	http.HandleFunc("/music/", mediaparty.Mp3Handler)
	http.HandleFunc("/table/", mediaparty.TableHandler)
	http.HandleFunc("/query/", mediaparty.QueryHandler)
	http.HandleFunc("/", mediaparty.MainHandler)
	http.ListenAndServe(":8000", nil)
}
