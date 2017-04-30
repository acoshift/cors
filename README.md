# cors

[![Go Report Card](https://goreportcard.com/badge/github.com/acoshift/cors)](https://goreportcard.com/report/github.com/acoshift/cors)
[![GoDoc](https://godoc.org/github.com/acoshift/cors?status.svg)](https://godoc.org/github.com/acoshift/cors)

CORS middleware for Golang net/http

### Example

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/acoshift/cors"
	"github.com/acoshift/middleware"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

func main() {
	h := middleware.Chain(
		cors.New(cors.Config{
			AllowOrigins:     []string{"localhost:8080"},
			AllowMethods:     []string{http.MethodGet, http.MethodPost},
			AllowHeaders:     []string{"Authorization"},
			AllowCredentials: true,
		}),
	)(http.HandlerFunc(handler))
	http.ListenAndServe(":8080", h)
}
```
