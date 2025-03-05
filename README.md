# SSE - Simple and Compliant Server-Sent Events Parser for Go
[![Go Report Card](https://goreportcard.com/badge/github.com/dblokhin/sse)](https://goreportcard.com/report/github.com/dblokhin/sse)
[![GoDoc](https://godoc.org/github.com/dblokhin/sse?status.svg)](https://godoc.org/github.com/dblokhin/sse)

**Keywords:** Golang SSE, Server-Sent Events, SSE Reader, SSE Parser, Go SSE Client

## Overview
**sse** is a lightweight, fully spec compliant Server-Sent Events (SSE) reader library for Go. It provides a simple wrapper over an `io.Reader` (typically an HTTP response body) so you can easily handle streaming SSE events in your Go applications. The design is really simple, easy to use, and perfect for developers looking to implement real-time event handling in a minimalistic way.

## Installation

```sh
go get github.com/dblokhin/sse
```

---

## Usage Example

`sse.Read` returns iterator, so it is a straightforward way to handle SSE stream data from an `io.Reader` source:

```go
	// ... 
	// Use the sse.Read iterator directly in a range loop.
	for event, err := range sse.Read(resp.Body) {
		if err != nil {
			// error handling ...
		}
		
		// process the valid event.
		fmt.Println("received event:", event)
	}
}
```

---

## Error Handling
Any invalid SSE line sequence is returned as an error (`ErrInvalidSequence`) with detailed message content, allowing graceful handling or logging of issues.

---

## License
MIT License# sse
