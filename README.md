Shell
=====

Simple C2 server example...

## Getting Started

Ensure `dep` is installed and use it to install the 3rd party dependencies
```
$ dep ensure
```

Run the server
```
$ cd cmd/server
$ go run *.go
```

Run a worker
```
$ cd cmd/shell
$ go run *.go
```

Run `ping` on the worker for 5 seconds
```
$ cd cmd/shell
$ go run *.go -timeout 5 -- ping google.com
```