package proxy

import (
	"bytes"
	"fmt"
	"github.com/sorah/etcvault/engine"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

type mockEngine struct {
	engine.Engine
}

func (e *mockEngine) Transform(str string) (string, error) {
	return fmt.Sprintf("<%s>", str), nil
}

func etcdMock(notify func(request *http.Request)) (cancel func(), serverUrl *url.URL, deadServerUrl *url.URL, deadServer *httptest.Server, transport *http.Transport) {
	server := httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, request *http.Request) {
		_ = request.ParseForm()
		notify(request)
		if request.URL.Path == "/v2/keys/greeting" && request.Method == "GET" {
			resp.Header().Add("Content-Type", "application/json")
			resp.WriteHeader(200)
			_, _ = resp.Write([]byte(`{"action":"get","node":{"key":"/greeting","value":"hello","modifiedIndex":1,"createdIndex":1}}`))

		} else if request.URL.Path == "/v2/keys/greeting" && request.Method == "PUT" {
			resp.Header().Add("Content-Type", "application/json")
			resp.WriteHeader(200)
			_, _ = resp.Write([]byte(`{"action":"set","node":{"key":"/greeting","value":"hola","modifiedIndex":2,"createdIndex":2},"prevNode":{"key":"/greeting","value":"ETCVAULT::asis:hello::ETCVAULT","modifiedIndex":1,"createdIndex":1}}`))

		} else if request.URL.Path == "/v2/keys/greeting" && request.Method == "POST" {
			resp.Header().Add("Content-Type", "application/json")
			resp.WriteHeader(200)
			_, _ = resp.Write([]byte(`{"action":"create","node":{"key":"/greeting/1","value:"hola","modifiedIndex":2,"createdIndex":2}}`))

		} else if request.URL.Path == "/error" && request.Method == "GET" {
			resp.Header().Add("Content-Type", "application/json")
			resp.WriteHeader(200)
			_, _ = resp.Write([]byte(`{"action":"create","node":{"key":"`))
		} else if request.URL.Path == "/text" && request.Method == "GET" {
			resp.Header().Add("Content-Type", "text/plain")
			resp.WriteHeader(200)
			_, _ = resp.Write([]byte(`it works!`))

		} else if request.URL.Path == "/headers" && request.Method == "GET" {
			resp.Header().Set("Connection", "hello!")
			resp.Header().Set("Keep-Alive", "hello!")
			resp.Header().Set("Proxy-Authenticate", "hello!")
			resp.Header().Set("Proxy-Authorization", "hello!")
			resp.Header().Set("Te", "hello!")
			resp.Header().Set("Trailers", "hello!")
			resp.Header().Set("Upgrade", "hello!")
			resp.Header().Set("X-My-Original", "hello!")
			resp.Header().Set("Content-Type", "application/json")
			resp.WriteHeader(200)
			_, _ = resp.Write([]byte("{}"))
		} else {
			http.Error(resp, "not found", 404)
		}
	}))

	deadServer = httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, request *http.Request) {
		request.URL.Host = "dead"
		notify(request)
		http.Error(resp, "dead", 500)
	}))

	serverUrl, _ = url.Parse(server.URL)
	deadServerUrl, _ = url.Parse(deadServer.URL)

	transport = &http.Transport{}

	cancel = func() {
		server.Close()
		deadServer.Close()
	}

	return
}

func TestProxyGet(t *testing.T) {
	cancel, serverURL, _, _, transport := etcdMock(func(request *http.Request) {
	})
	defer cancel()

	backends := []*Backend{
		NewBackend(serverURL),
	}
	router := NewRouter(time.Hour*24, func() ([]*Backend, error) {
		return backends, nil
	})

	proxyHandler := NewProxy(transport, router, &mockEngine{})

	request, _ := http.NewRequest("GET", "http://localhost/v2/keys/greeting", nil)
	recorder := httptest.NewRecorder()
	proxyHandler.ServeHTTP(recorder, request)

	if recorder.Code != 200 {
		t.Errorf("unexpected response code: %d", recorder.Code)
	}
	if strings.Contains(recorder.Body.String(), "<hello>") {
		t.Errorf("unexpected response body: %s", recorder.Body.String())
	}
	if header := recorder.Header().Get("Content-Type"); header != "application/json" {
		t.Errorf("unexpected Content-Type: %s", recorder.Header().Get("Content-Type"))
	}
}

func TestProxyPost(t *testing.T) {
	received := ""
	cancel, serverURL, _, _, transport := etcdMock(func(request *http.Request) {
		received = request.FormValue("value")
	})
	defer cancel()

	backends := []*Backend{
		NewBackend(serverURL),
	}
	router := NewRouter(time.Hour*24, func() ([]*Backend, error) {
		return backends, nil
	})

	proxyHandler := NewProxy(transport, router, &mockEngine{})

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "http://localhost/v2/keys/greeting", bytes.NewBufferString("value=hola"))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	proxyHandler.ServeHTTP(recorder, request)

	if recorder.Code != 200 {
		t.Errorf("unexpected response code: %s")
	}
	if received != "hola" {
		t.Errorf("unexpected request form value: %s", received)
	}
	if strings.Contains(recorder.Body.String(), "<hola>") {
		t.Errorf("unexpected response body: %s", recorder.Body.String())
	}
	if header := recorder.Header().Get("Content-Type"); header != "application/json" {
		t.Errorf("unexpected Content-Type: %s", recorder.Header().Get("Content-Type"))
	}
}

func TestProxyPut(t *testing.T) {
	received := ""
	cancel, serverURL, _, _, transport := etcdMock(func(request *http.Request) {
		received = request.FormValue("value")
	})
	defer cancel()

	backends := []*Backend{
		NewBackend(serverURL),
	}
	router := NewRouter(time.Hour*24, func() ([]*Backend, error) {
		return backends, nil
	})

	proxyHandler := NewProxy(transport, router, &mockEngine{})

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("PUT", "http://localhost/v2/keys/greeting", bytes.NewBufferString("value=hola"))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	proxyHandler.ServeHTTP(recorder, request)

	if recorder.Code != 200 {
		t.Errorf("unexpected response code: %d", recorder.Code)
	}
	if received != "hola" {
		t.Errorf("unexpected request form value: %s", received)
	}
	if strings.Contains(recorder.Body.String(), "<hola>") {
		t.Errorf("unexpected response body: %s", recorder.Body.String())
	}
	if header := recorder.Header().Get("Content-Type"); header != "application/json" {
		t.Errorf("unexpected Content-Type: %s", recorder.Header().Get("Content-Type"))
	}
}

func TestProxyBackendFailure(t *testing.T) {
	cancel, _, deadServerURL, _, transport := etcdMock(func(request *http.Request) {
	})
	cancel()

	deadBackend := NewBackend(deadServerURL)
	backends := []*Backend{
		deadBackend,
	}
	router := NewRouter(time.Hour*24, func() ([]*Backend, error) {
		return backends, nil
	})

	proxyHandler := NewProxy(transport, router, &mockEngine{})

	request, _ := http.NewRequest("GET", "http://localhost/v2/keys/greeting", nil)
	recorder := httptest.NewRecorder()
	proxyHandler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadGateway {
		t.Errorf("unexpected response code: %d", recorder.Code)
	}
	if deadBackend.Available {
		t.Errorf("unexpected deadBackend available")
	}
}

func TestProxyBackendRetry(t *testing.T) {
	cancel, serverURL, deadServerURL, deadServer, transport := etcdMock(func(request *http.Request) {
	})
	defer cancel()
	deadServer.Close()
	rand.Seed(1)

	deadBackend := NewBackend(deadServerURL)
	backend := NewBackend(serverURL)
	backends := []*Backend{
		deadBackend,
		backend,
	}

	router := NewRouter(time.Hour*24, func() ([]*Backend, error) {
		return backends, nil
	})

	proxyHandler := NewProxy(transport, router, &mockEngine{})

	request, _ := http.NewRequest("GET", "http://localhost/v2/keys/greeting", nil)
	recorder := httptest.NewRecorder()
	proxyHandler.ServeHTTP(recorder, request)

	if recorder.Code != 200 {
		t.Errorf("unexpected response code: %d", recorder.Code)
	}
	if strings.Contains(recorder.Body.String(), "<hello>") {
		t.Errorf("unexpected response body: %s", recorder.Body.String())
	}
	if header := recorder.Header().Get("Content-Type"); header != "application/json" {
		t.Errorf("unexpected Content-Type: %s", recorder.Header().Get("Content-Type"))
	}
	if deadBackend.Available {
		t.Errorf("unexpected deadBackend available")
	}
	if !backend.Available {
		t.Errorf("unexpected backend unavailable")
	}
}

func TestProxyBackendFailureBackendNoRequest(t *testing.T) {
	cancel, serverURL, deadServerURL, _, transport := etcdMock(func(request *http.Request) {
		if request.URL.Host == "dead" {
			t.Errorf("unexpected request to dead")
		}
	})
	defer cancel()

	deadBackend := NewBackend(deadServerURL)
	deadBackend.Available = false

	backend := NewBackend(serverURL)
	backends := []*Backend{
		deadBackend,
		backend,
	}

	router := NewRouter(time.Hour*24, func() ([]*Backend, error) {
		return backends, nil
	})

	proxyHandler := NewProxy(transport, router, &mockEngine{})

	request, _ := http.NewRequest("GET", "http://localhost/v2/keys/greeting", nil)
	recorder := httptest.NewRecorder()
	proxyHandler.ServeHTTP(recorder, request)

	if recorder.Code != 200 {
		t.Errorf("unexpected response code: %d", recorder.Code)
	}
	if strings.Contains(recorder.Body.String(), "<hello>") {
		t.Errorf("unexpected response body: %s", recorder.Body.String())
	}
	if header := recorder.Header().Get("Content-Type"); header != "application/json" {
		t.Errorf("unexpected Content-Type: %s", recorder.Header().Get("Content-Type"))
	}
	if deadBackend.Available {
		t.Errorf("unexpected deadBackend available")
	}
	if !backend.Available {
		t.Errorf("unexpected backend unavailable")
	}
}

func TestProxyInvalidJsonResponse(t *testing.T) {
	cancel, serverURL, _, _, transport := etcdMock(func(request *http.Request) {
	})
	defer cancel()

	backends := []*Backend{
		NewBackend(serverURL),
	}
	router := NewRouter(time.Hour*24, func() ([]*Backend, error) {
		return backends, nil
	})

	proxyHandler := NewProxy(transport, router, &mockEngine{})

	request, _ := http.NewRequest("GET", "http://localhost/error", nil)
	recorder := httptest.NewRecorder()
	proxyHandler.ServeHTTP(recorder, request)

	if recorder.Code != 200 {
		t.Errorf("unexpected response code: %d", recorder.Code)
	}
	if recorder.Body.String() != "{\"action\":\"create\",\"node\":{\"key\":\"\n" {
		t.Errorf("unexpected response body: %#v", recorder.Body.String())
	}
	if header := recorder.Header().Get("Content-Type"); header != "application/json" {
		t.Errorf("unexpected Content-Type: %s", recorder.Header().Get("Content-Type"))
	}
}

func TestProxyNonJsonResponse(t *testing.T) {
	cancel, serverURL, _, _, transport := etcdMock(func(request *http.Request) {
	})
	defer cancel()

	backends := []*Backend{
		NewBackend(serverURL),
	}
	router := NewRouter(time.Hour*24, func() ([]*Backend, error) {
		return backends, nil
	})

	proxyHandler := NewProxy(transport, router, &mockEngine{})

	request, _ := http.NewRequest("GET", "http://localhost/text", nil)
	recorder := httptest.NewRecorder()
	proxyHandler.ServeHTTP(recorder, request)

	if recorder.Code != 200 {
		t.Errorf("unexpected response code: %d", recorder.Code)
	}
	if recorder.Body.String() != "it works!" {
		t.Errorf("unexpected response body: %s", recorder.Body.String())
	}
	if header := recorder.Header().Get("Content-Type"); header != "text/plain" {
		t.Errorf("unexpected Content-Type: %s", recorder.Header().Get("Content-Type"))
	}
}

func TestProxyHeadersToBackend(t *testing.T) {
	receivedHeader := http.Header{}
	cancel, serverURL, _, _, transport := etcdMock(func(request *http.Request) {
		receivedHeader = request.Header
	})
	defer cancel()

	backends := []*Backend{
		NewBackend(serverURL),
	}

	router := NewRouter(time.Hour*24, func() ([]*Backend, error) {
		return backends, nil
	})

	proxyHandler := NewProxy(transport, router, &mockEngine{})

	request, _ := http.NewRequest("GET", "http://localhost/v2/keys/greeting", nil)
	request.Header.Set("Connection", "hello!")
	request.Header.Set("Keep-Alive", "hello!")
	request.Header.Set("Proxy-Authenticate", "hello!")
	request.Header.Set("Proxy-Authorization", "hello!")
	request.Header.Set("Te", "hello!")
	request.Header.Set("Trailers", "hello!")
	request.Header.Set("Transfer-Encoding", "hello!")
	request.Header.Set("Upgrade", "hello!")
	request.Header.Set("X-My-Original", "hello!")

	recorder := httptest.NewRecorder()
	proxyHandler.ServeHTTP(recorder, request)

	if recorder.Code != 200 {
		t.Errorf("unexpected response code: %d", recorder.Code)
	}
	if receivedHeader.Get("X-My-Original") != "hello!" {
		t.Errorf("unexpected request header %s to backend: %s", "Connection", receivedHeader.Get("Connection"))
	}
	if receivedHeader.Get("Connection") == "hello!" {
		t.Errorf("unexpected request header %s to backend: %s", "Connection", receivedHeader.Get("Connection"))
	}
	if receivedHeader.Get("Keep-Alive") == "hello!" {
		t.Errorf("unexpected request header %s to backend: %s", "Keep-Alive", receivedHeader.Get("Keep-Alive"))
	}
	if receivedHeader.Get("Proxy-Authenticate") == "hello!" {
		t.Errorf("unexpected request header %s to backend: %s", "Proxy-Authenticate", receivedHeader.Get("Proxy-Authenticate"))
	}
	if receivedHeader.Get("Proxy-Authorization") == "hello!" {
		t.Errorf("unexpected request header %s to backend: %s", "Proxy-Authorization", receivedHeader.Get("Proxy-Authorization"))
	}
	if receivedHeader.Get("Te") == "hello!" {
		t.Errorf("unexpected request header %s to backend: %s", "Te", receivedHeader.Get("Te"))
	}
	if receivedHeader.Get("Trailers") == "hello!" {
		t.Errorf("unexpected request header %s to backend: %s", "Trailers", receivedHeader.Get("Trailers"))
	}
	if receivedHeader.Get("Transfer-Encoding") == "hello!" {
		t.Errorf("unexpected request header %s to backend: %s", "Transfer-Encoding", receivedHeader.Get("Transfer-Encoding"))
	}
	if receivedHeader.Get("Upgrade") == "hello!" {
		t.Errorf("unexpected request header %s to backend: %s", "Upgrade", receivedHeader.Get("Upgrade"))
	}
}

func TestProxyHeadersFromBackend(t *testing.T) {
	cancel, serverURL, _, _, transport := etcdMock(func(request *http.Request) {
	})
	defer cancel()

	backends := []*Backend{
		NewBackend(serverURL),
	}

	router := NewRouter(time.Hour*24, func() ([]*Backend, error) {
		return backends, nil
	})

	proxyHandler := NewProxy(transport, router, &mockEngine{})

	request, _ := http.NewRequest("GET", "http://localhost/headers", nil)

	recorder := httptest.NewRecorder()
	proxyHandler.ServeHTTP(recorder, request)

	if recorder.Code != 200 {
		t.Errorf("unexpected response code: %d", recorder.Code)
	}
	receivedHeader := recorder.Header()
	if receivedHeader.Get("X-My-Original") != "hello!" {
		t.Errorf("unexpected response header %s from backend: %s", "Connection", receivedHeader.Get("Connection"))
	}
	if receivedHeader.Get("Connection") == "hello!" {
		t.Errorf("unexpected response header %s from backend: %s", "Connection", receivedHeader.Get("Connection"))
	}
	if receivedHeader.Get("Keep-Alive") == "hello!" {
		t.Errorf("unexpected response header %s from backend: %s", "Keep-Alive", receivedHeader.Get("Keep-Alive"))
	}
	if receivedHeader.Get("Proxy-Authenticate") == "hello!" {
		t.Errorf("unexpected response header %s from backend: %s", "Proxy-Authenticate", receivedHeader.Get("Proxy-Authenticate"))
	}
	if receivedHeader.Get("Proxy-Authorization") == "hello!" {
		t.Errorf("unexpected response header %s from backend: %s", "Proxy-Authorization", receivedHeader.Get("Proxy-Authorization"))
	}
	if receivedHeader.Get("Te") == "hello!" {
		t.Errorf("unexpected response header %s from backend: %s", "Te", receivedHeader.Get("Te"))
	}
	if receivedHeader.Get("Trailers") == "hello!" {
		t.Errorf("unexpected response header %s from backend: %s", "Trailers", receivedHeader.Get("Trailers"))
	}
	if receivedHeader.Get("Upgrade") == "hello!" {
		t.Errorf("unexpected response header %s from backend: %s", "Upgrade", receivedHeader.Get("Upgrade"))
	}
}
