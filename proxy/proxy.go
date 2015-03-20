package proxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sorah/etcvault/engine"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type ClosableBuffer struct {
	*bytes.Buffer
}

func (buf ClosableBuffer) Close() error {
	return nil
}

// Hop-by-hop headers (borrowed from httputil.ReverseProxy)
// http://www.w3.org/Protocols/rfc2616/rfc2616-sec13.html
var singleHopHeaders = []string{
	"Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te",
	"Trailers",
	"Transfer-Encoding",
	"Upgrade",
}

type Proxy struct {
	Transport    *http.Transport
	Router       *Router
	Engine       engine.Transformable
	AdvertiseUrl string
}

func NewProxy(transport *http.Transport, router *Router, e engine.Transformable, advertiseUrl string) http.Handler {
	return &Proxy{
		Transport:    transport,
		Router:       router,
		Engine:       e,
		AdvertiseUrl: advertiseUrl,
	}
}

func (proxy *Proxy) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	if request.URL.Path == "/v2/members" {
		proxy.serveMembersRequest(response, request)
	} else if request.URL.Path == "/v2/machines" {
		proxy.serveMachinesRequest(response, request)
	} else if request.URL.Path == "/_etcvault/keys" {
		proxy.serveEtcvaultKeysRequest(response, request)
	} else {
		proxy.serveProxyRequest(response, request)
	}
}

func (proxy *Proxy) serveProxyRequest(response http.ResponseWriter, request *http.Request) {
	backendRequest := new(http.Request)
	// copy
	*backendRequest = *request
	backendRequest.Header = make(http.Header)

	backendRequest.Proto = "HTTP/1.1"
	backendRequest.ProtoMajor = 1
	backendRequest.ProtoMinor = 1
	backendRequest.Close = false

	copyHeader(request.Header, backendRequest.Header)
	removeSingleHopHeaders(&backendRequest.Header)

	if (backendRequest.Method == "POST" || backendRequest.Method == "PUT" || backendRequest.Method == "PATCH") && backendRequest.Body != nil {
		origBody := backendRequest.Body
		defer origBody.Close()

		if err := backendRequest.ParseForm(); err != nil {
			log.Printf("couldn't parse form: %s", err.Error())
			http.Error(response, "couldn't parse form", 400)
			return
		}

		if backendRequest.PostForm != nil {
			origValue := backendRequest.PostForm.Get("value")
			value, err := proxy.Engine.Transform(origValue)
			if err == nil {
				backendRequest.PostForm.Set("value", value)
			} else {
				log.Printf("failed to transform value: %s", err.Error())
			}
			newFormString := backendRequest.PostForm.Encode()
			backendRequest.Body = ClosableBuffer{bytes.NewBufferString(newFormString)}
			backendRequest.ContentLength = int64(len(newFormString))
		}
	}

	var backendResponse *http.Response

	backends := proxy.Router.ShuffledAvailableBackends()
	for _, backend := range backends {
		backendRequest.URL.Scheme = backend.Url.Scheme
		backendRequest.URL.Host = backend.Url.Host

		var err error
		backendResponse, err = proxy.Transport.RoundTrip(backendRequest)
		if err != nil {
			log.Printf("backend %s response error: %s", backend.Url.String(), err.Error())
			backend.Fail()
			continue
		}
		backend.Ok()
		break
	}

	if backendResponse == nil {
		log.Printf("all backends not available...")
		http.Error(response, "backends all unavailable", http.StatusBadGateway)
		return
	}

	defer backendResponse.Body.Close()

	removeSingleHopHeaders(&backendResponse.Header)
	copyHeader(backendResponse.Header, response.Header())

	if backendResponse.Header.Get("Content-Type") == "application/json" {
		json, err := ioutil.ReadAll(backendResponse.Body)
		if err != nil {
			panic(err)
		}

		transformedJson, err := proxy.Engine.TransformEtcdJsonResponse(json)
		if err == nil {
			response.Header().Set("Content-Length", fmt.Sprintf("%d", len(transformedJson)+1))
			response.WriteHeader(backendResponse.StatusCode)
			response.Write(transformedJson)
			response.Write([]byte("\n"))
		} else {
			fmt.Printf("transform error %s\n", err.Error())
			response.WriteHeader(backendResponse.StatusCode)
			response.Write(json)
		}
	} else {
		response.WriteHeader(backendResponse.StatusCode)
		io.Copy(response, backendResponse.Body)
	}
}

func (proxy *Proxy) serveMembersRequest(response http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		http.Error(response, "not supported; communicate with etcd directly", http.StatusMethodNotAllowed)
		return
	}

	type memberT struct {
		ClientURLs []string
		PeerURLs   []string
		Name       string
		Id         string
	}

	jsonBytes, err := json.Marshal(struct {
		Members []memberT
	}{
		Members: []memberT{
			{
				ClientURLs: []string{proxy.AdvertiseUrl},
				Name:       "etcvault",
				Id:         "deadbeef",
			},
		},
	})

	if err != nil {
		http.Error(response, "failed to marshal", 500)
		log.Printf("failed to marshal /v2/members: %s", err.Error())
		return
	}

	response.Header().Add("Content-Type", "application/json")
	response.Header().Add("Server", "etcvault")
	response.WriteHeader(200)
	response.Write(jsonBytes)
}

func (proxy *Proxy) serveMachinesRequest(response http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		http.Error(response, "not supported; communicate with etcd directly", http.StatusMethodNotAllowed)
		return
	}

	response.Header().Add("Content-Type", "text/plain")
	response.Header().Add("Server", "etcvault")
	response.WriteHeader(200)
	response.Write([]byte(proxy.AdvertiseUrl))
}

func (proxy *Proxy) serveEtcvaultKeysRequest(response http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		http.Error(response, "not found", http.StatusNotFound)
		return
	}

	request.ParseForm()
	var list []string
	if request.FormValue("encryption") != "" {
		list = proxy.Engine.GetKeychain().ListForEncryption()
	} else {
		list = proxy.Engine.GetKeychain().List()
	}

	response.Header().Add("Content-Type", "text/plain")
	response.Header().Add("Server", "etcvault")
	response.WriteHeader(200)

	for _, name := range list {
		response.Write([]byte(name))
		response.Write([]byte("\n"))
	}
}

func copyHeader(source, destination http.Header) {
	for key, values := range source {
		for _, value := range values {
			destination.Add(key, value)
		}
	}
}

func removeSingleHopHeaders(header *http.Header) {
	for _, name := range singleHopHeaders {
		header.Del(name)
	}
}
