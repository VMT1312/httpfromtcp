package main

import (
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const (
	port       = 42069
	badRequest = `<html>
  <head>
	<title>400 Bad Request</title>
  </head>
  <body>
	<h1>Bad Request</h1>
	<p>Your request honestly kinda sucked.</p>
  </body>
</html>
`
	internalError = `<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>
	`
	goodRequest = `<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>
	`
)

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handler(w *response.Writer, req *request.Request) {
	if req.RequestLine.RequestTarget == "/yourproblem" {
		w.WriteStatusLine(response.StatusBadRequest)
		h := response.GetDefaultHeaders(len(badRequest))
		w.WriteHeaders(h)
		w.WriteBody([]byte(badRequest))
		return
	}
	if req.RequestLine.RequestTarget == "/myproblem" {
		w.WriteStatusLine(response.StatusInternalError)
		h := response.GetDefaultHeaders(len(internalError))
		w.WriteHeaders(h)
		w.WriteBody([]byte(internalError))
		return
	}
	w.WriteStatusLine(response.StatusOK)
	h := response.GetDefaultHeaders(len(goodRequest))
	w.WriteHeaders(h)
	w.WriteBody([]byte(goodRequest))
}
