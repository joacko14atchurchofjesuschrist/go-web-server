# Simple Go Web Server

A lightweight web server implemented in Go that serves static files and responds to HTTP requests.

## Project Overview

This project demonstrates a basic web server built with Go's standard library. The server has the following features:

- Responds to requests at the root path with a "Hello World" message
- Serves static files from the `static/` directory
- Listens for connections on port 80

## Project Structure

```
.
├── go.mod              # Go module definition
├── README.md           # This file
├── server.go           # Main server implementation
└── static/             # Directory for static files
    └── index.html      # Sample HTML file served by the server
```

## Getting Started

### Prerequisites

- Go installed on your system (download from [https://golang.org/dl/](https://golang.org/dl/))

### Running the Server

To run the web server, execute:

```bash
go run server.go
```

The server will start and listen on port 80. You can access it by visiting:

- [http://localhost/](http://localhost/) - Shows a "Hello World" message along with the requested path
- [http://localhost/static/index.html](http://localhost/static/index.html) - Shows the content of the static HTML file

## Modifying the Server

- To add more static files, simply place them in the `static/` directory
- To add more route handlers, add additional `http.HandleFunc()` calls in `server.go`

## Note

Since the server runs on port 80, you might need administrator privileges to start it. If you encounter permission issues, you can modify the port number in `server.go` to use a higher port number (e.g., 8080).