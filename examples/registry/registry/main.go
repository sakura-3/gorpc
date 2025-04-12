package main

import (
	"gorpc/registry"
	"net/http"
)

func main() {
	// serve at http://localhost:9999/_gorpc_/registry
	if err := http.ListenAndServe(":9999", registry.DefaultGeeRegister); err != nil {
		panic(err)
	}
}
