package main

import (
	"code.google.com/p/go-uuid/uuid"
	"fmt"
	"net/http"
)

var lookup map[string]string = make(map[string]string)

func handler(w http.ResponseWriter, r *http.Request) {
	longUrl := r.URL.Path[9:]
	//Do some validation on the long url
	shortUrl := shortenUrl(longUrl)

	fmt.Fprintf(w, "shortened Url", shortUrl)
}

func lookupHandler(w http.ResponseWriter, r *http.Request) {
	shortUrl := r.URL.Path[8:]
	//TODO: Validate that the short url is non null
	longUrl := lookItUp(shortUrl)
	fmt.Fprintf(w, "output: ", longUrl)
}

func shortenUrl(longUrl string) string {
	shortUrl := createUniqueMapping(longUrl)
	lookup[shortUrl] = longUrl
	return shortUrl
}

func createUniqueMapping(longUrl string) string {
	id := uuid.NewUUID()
	return id.String()[:7]
}
func lookItUp(shortUrl string) string {
	if longUrl := lookup[shortUrl]; longUrl != "" {
		return longUrl
	} else {
		return "unable to lookup: " + shortUrl
	}
}

func main() {
	http.HandleFunc("/shorten/", handler)
	http.HandleFunc("/lookup/", lookupHandler)
	http.ListenAndServe(":8800", nil)
}
