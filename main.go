package main

import (
	"code.google.com/p/go-uuid/uuid"
	"fmt"
	"github.com/mediocregopher/radix.v2/redis"
	"net/http"
)

var lookup map[string]string = make(map[string]string)
var client *redis.Client

const store = "redis"

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
	for lookupIsNonEmpty(shortUrl) {
		shortUrl = createUniqueMapping(longUrl)
	}
	addMappingToStore(shortUrl, longUrl)
	return shortUrl
}

func addMappingToStore(shortUrl, longUrl string) {
	switch store {
	case "redis":
		err := client.Cmd("Set", shortUrl, longUrl).Err
		if err != nil {
			fmt.Printf("error while writing to redis: " + err.Error())
		}
	default:
		lookup[shortUrl] = longUrl
	}
}

func lookupIsNonEmpty(shortUrl string) bool {
	switch store {
	case "redis":
		longUrl, _ := client.Cmd("GET", shortUrl).Str()
		return longUrl != ""
	default:
		return lookup[shortUrl] != ""
	}
}

func createUniqueMapping(longUrl string) string {
	id := uuid.NewUUID()
	return id.String()[:7]
}

func lookItUp(shortUrl string) string {
	//Later on one can have any number of schemes for generating this shortened url
	switch store {
	case "redis":
		longUrl, err := client.Cmd("GET", shortUrl).Str()
		if err != nil {
			return "unable to lookup: " + shortUrl
		}
		return longUrl
	default:
		if longUrl := lookup[shortUrl]; longUrl != "" {
			return longUrl
		} else {
			return "unable to lookup: " + shortUrl
		}
	}
}
func setupRedis() {

	c, err := redis.Dial("tcp", "localhost:6379")
	client = c
	if err != nil {
		//handle error
	}
}

func main() {
	setupRedis()
	http.HandleFunc("/shorten/", handler)
	http.HandleFunc("/lookup/", lookupHandler)
	http.ListenAndServe(":8800", nil)
}
