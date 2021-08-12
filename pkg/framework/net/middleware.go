package net

import (
	"compress/gzip"
	"net/http"
)

func DecompressMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.Header.Get("Content-Encoding") {
		case "gzip":
			gz, err := gzip.NewReader(request.Body)
			if err != nil {
				writer.WriteHeader(http.StatusBadRequest)
			}
			defer gz.Close()
			request.Body = gz
			break
		}
		next.ServeHTTP(writer, request)
	})
}

func JsonResponseTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(writer, request)
	})
}