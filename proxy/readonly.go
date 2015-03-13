package proxy

import (
	"github.com/sorah/etcvault/engine"
	"net/http"
)

func NewReadonlyProxy(transport *http.Transport, router *Router, e engine.Transformable, advertiseUrl string) http.Handler {
	return readonlyHandler(NewProxy(transport, router, e, advertiseUrl))
}

func readonlyHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.Method != "GET" {
			// I prefer method not allowed, but following etcd's proxy mode behavior for compat
			response.WriteHeader(http.StatusNotImplemented)
			return
		}

		handler.ServeHTTP(response, request)
	})
}
