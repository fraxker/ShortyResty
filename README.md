# ShortyResty

ShortyResty is a url shortener written in Golang.

## Endpoints
ShortyResty hosts two endpoints, one for shortening the url and one for going to the shortened url.

### Shorten `/shorten`
This endpoint takes in a POST request with the json format:
`{“url”: “http://example.com/verylonguselessURLthatdoesnotseemtoend/123”}`

The response will be formated as such:
`{“short_url”: “http://127.0.0.1:8080/xxxxxxxx”}`
where `xxxxxxxx` is the id number generated.

### Retrieve `/{ID}`
The retrieve endpoint takes in a GET request and will 302 redirect you to the url associated with that id.

## Usage
ShortyResty can be run via:
```sh
go run .
```
And will be accessible at http://127.0.0.1:8080

## Dependencies

ShortyResty depends on standard libraries with exception of gorilla/mux. Dependencies can be installed via
```sh
go get
```

## Future Improvements
Potential improvements include:
- Backing with database to be persistent across reboots
- Creating home page/other web pages
- Adding addition endpoint to check all current ids/urls
