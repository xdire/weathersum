# Weathersum
An example of simple weather summarization tool

# Release stages
1) First release with `/v1/weather` API

2) Further release with simple LLM summarization 
as a bundled feature

# How to run
### Build the app
```
go build .
```
### Run the app
```
./weathersum
```
### Run main function with go run
```
go run main.go
```

# How to run tests
```
go test ./...
```
> Functional tests included in the [tests](/tests) directory


# How to interact

After running the application the API should be reachable as following:
```
http://localhost:8000/v1/weather?lat=37.7739&lon=-122.4313
```
Or with curl
```
curl 'http://localhost:8000/v1/weather?lat=37.7739&lon=-122.4313'
```
