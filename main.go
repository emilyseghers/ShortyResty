package main

import (
	"fmt"
	"encoding/json"
	"time"
	"log"
	"math/rand"
	"net/url"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/golang/gddo/httputil/header"
)


//Map variable that will keep track of shortened URLs by ID.
//Keys of the map will by the IDs and values will be the 
//original long URLs.
var paths = make(map[string]string)


//Charset that will be used to generate random shortened 
const letters = "ABCDEFGHIJKLMOPQRSTUVWXYZ"


//Handles incoming GET requests of the format "http://127.0.0.1:8080/ID".
//If the given ID maps to a long URL this method will 302 redirect the user
//to the original URL, otherwise it will throw an error.
func redirect(w http.ResponseWriter, r *http.Request) {
	//Retrieves the the map 
	vars := mux.Vars(r)
	id := vars["id"]

	//If the ID from the GET request is pressent in the map it
	//will redirect the user to the long URL link and otherwise it will
	//throw an error.
	if val, ok := paths[id]; ok {
		http.Redirect(w, r, val, http.StatusFound)
		return
	}
	http.Error(w, "ID didn't map to a URL", http.StatusBadRequest)
	return
}


//Handles incoming POST requests in JSON format.
//After recieving this request it will response with a 
//shortened URL also in JSON format that is created using the 
//makeID() method below. Keeps track of the inputted URL and 
//outputted shortened URL by utilizing a map structure.
func shorten(w http.ResponseWriter, r *http.Request) {
	var id string
	var request Req
	var jurl Resp
	var short_url string

	//Error checking that the Content-Type header matches that of 
	//a JSON file.
	if r.Header.Get("Content-Type") != "" {
        value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
        if value != "application/json" {
            msg := "Content-Type header is not application/json"
            http.Error(w, msg, http.StatusUnsupportedMediaType)
            return
        }
    }

	//Error checking that the JSON request file is not too large
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	//Decodes the JSON request
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//Checking that the URL is valid
	_, err = url.ParseRequestURI(request.Url)
	if err != nil {
		http.Error(w, "URL is not valid", http.StatusBadRequest)
		return
	}

	//Generates a randomized ID and if the produced ID already 
	//exists it creates a new one
	id = makeID(8)
	for key := range paths {
		if ( key == id){
			id = makeID(8)
		}
	}

	//Creates the shortened URL using the ID and adds it to the map
	short_url = "http://127.0.0.1:8080/" + id
	paths[id] = request.Url

	//Creates the JSON response and sends it
	jurl.Short_url = short_url
	jData, err := json.Marshal(jurl)
	if err != nil {
		panic(err)
		return
	}
	w.Write(jData)
}


//Generates an 8 character randomized ID composed of all 
//uppercase letters of the alphabet.
func makeID(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}


//Registers two routers mapping URL paths to handlers
//using the gorilla/mux package.
func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/shorten", shorten).Methods("POST")
	router.HandleFunc("/{id}", redirect).Methods("GET")
	err := http.ListenAndServe(":8080", router)
	log.Fatal(err)
}


//Main method which lets the user know what port the web
//server is available on and calls a method which handles 
//the incoming requests.
func main() {
	fmt.Printf("Starting server on port 8080")
	handleRequests()
}


//Struct for http Requests, only content is the given URL.
type Req struct {
	Url string
}

//Struct for http Response, only content is the randomized 
//shortened URL.
type Resp struct {
	Short_url string
}
