# 1&1 Cloudserver API Go Client
  
This project contains a Go implementation of a client for the 1&1 Cloudserver's public API.
    
## About 1&1 Cloudserver API

The Cloudserver API Documentation can be found here: []()
      
The Go Client Documentation can be found here: []()

## Usage

Import the library:

```go
import oaocs "github.com/jlusiardi/oneandone-cloudserver-api"
```

Create a new API instance:
```go
api := oaocs.New("YOUR_TOKEN", "API_ENDPOINT")
```

Query all server:
```go
servers, err := api.GetServers()
if err != nil {
	// some error handling
}
```
  
## Contributing
    
We are happy to recieve your reports on any issue and also pull requests are welcome!
