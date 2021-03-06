# ShortyResty URL Shortener

Topics: JSON, httpHandlers, maps

## Overview

This program is a Golang REST API service that allows a user to shorten a URL. It also has the capability to redirect that user to the long URL after inputting the previously generated shortened URL.

## Specifications

This program is available on port 8080. It responds to a POST request with endpoint "/shorten." The request should be formatted in JSON otherwise an error will be thrown. The respose will also be in JSON and will contain a shortened URL that is randomly generated by the program and stored in a map.

The other feature of this program is redirection. If a user makes a GET request with endpoint "/$ID" (ID being the randomaly generated ID) of the format "http://127.0.0.1:8080/ID" then the program will 302 redirect that user to the long form of the URL.

## Request/Response Formatting

POST JSON request should be of the format: {"url": "http://example.com/verylonguselessURLthatdoesnotseemtoend/123"}

JSON response will be of the format: {"short_url": "http://127.0.0.1:8080/xxxxxxxx"} where xxxxxxxx is the ID

GET request to your browser of client: http://127.0.0.1:8080/ID 

This will redirect you to the long form of the URL if the ID is valid.

## Testing and Requirements 

This program was tested using Postman to generate POST and GET requests.

This program utilizes the gorilla/mux and golang/gddo/httputil/header github repositories.

Execute: 

go get -u github.com/gorilla/mux

go get -u github.com/golang/gddo/httputil/header 

in your terminal before running.
