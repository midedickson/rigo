package rigo

import (
	"net/http"
)

func handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		initialize()
	})
}
