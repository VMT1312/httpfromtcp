package main

import (
	"crypto/sha256"
	"fmt"
	"httpfromtcp/internal/headers"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
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
	server, err := server.Serve(port, handlerRequest)
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

func handlerRequest(w *response.Writer, req *request.Request) {
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/") {
		handlerProxy(w, req)
		return
	}
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
	if req.RequestLine.RequestTarget == "/video" {
		b, err := os.ReadFile("assets/vim.mp4")
		if err != nil {
			w.WriteStatusLine(response.StatusInternalError)
			h := response.GetDefaultHeaders(len(internalError))
			w.WriteHeaders(h)
			w.WriteBody([]byte(internalError))
			return
		}
		w.WriteStatusLine(response.StatusOK)
		h := response.GetDefaultHeaders(len(b))
		h.Override("Content-Type", "video/mp4")
		w.WriteHeaders(h)
		w.WriteBody(b)
		return
	}
	w.WriteStatusLine(response.StatusOK)
	h := response.GetDefaultHeaders(len(goodRequest))
	w.WriteHeaders(h)
	w.WriteBody([]byte(goodRequest))
}

func handlerProxy(w *response.Writer, req *request.Request) {
	subPath := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")
	url := "https://httpbin.org/" + subPath
	resp, err := http.Get(url)
	if err != nil {
		handlerRequest(w, &request.Request{
			RequestLine: request.RequestLine{
				RequestTarget: "/myproblem",
			},
		})
	}
	defer resp.Body.Close()
	buf := make([]byte, 1024)
	w.WriteStatusLine(response.StatusOK)
	h := response.GetDefaultHeaders(0)
	h.Remove("Content-Length")
	h.Override("Transfer-Encoding", "chunked")
	h.Override("Trailer", "X-Content-SHA256, X-Content-Length")
	w.WriteHeaders(h)
	fullBody := []byte{}
	for {
		n, err := resp.Body.Read(buf)
		if err != nil && err != io.EOF {
			handlerRequest(w, &request.Request{
				RequestLine: request.RequestLine{
					RequestTarget: "/myproblem",
				},
			})
			return
		}
		if n > 0 {
			fullBody = append(fullBody, buf[:n]...)
			w.WriteChunkedBody(buf[:n])
		}
		if err == io.EOF {
			w.WriteChunkedBodyDone()
			hashString := fmt.Sprintf("%x", sha256.Sum256(fullBody))
			contentLength := len(fullBody)
			w.WriteTrailers(headers.Headers{
				"X-Content-SHA256": hashString,
				"X-Content-Length": fmt.Sprintf("%d", contentLength),
			})
			return
		}
	}
}
