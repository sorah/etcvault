package proxy

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestReadonlyProxyGet(t *testing.T) {
	cancel, serverURL, _, _, transport := etcdMock(func(request *http.Request) {
		if request.Method != "GET" {
			t.Errorf("Received non GET request")
		}
	})
	defer cancel()

	backends := []*Backend{
		NewBackend(serverURL),
	}

	router := NewRouter(time.Hour*24, func() ([]*Backend, error) {
		return backends, nil
	})

	proxyHandler := NewReadonlyProxy(transport, router, &mockEngine{}, "http://localhost")

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

func TestReadonlyProxyPost(t *testing.T) {
	cancel, serverURL, _, _, transport := etcdMock(func(request *http.Request) {
		if request.Method != "GET" {
			t.Errorf("Received non GET request")
		}
	})
	defer cancel()

	backends := []*Backend{
		NewBackend(serverURL),
	}

	router := NewRouter(time.Hour*24, func() ([]*Backend, error) {
		return backends, nil
	})

	proxyHandler := NewReadonlyProxy(transport, router, &mockEngine{}, "http://localhost")

	request, _ := http.NewRequest("POST", "http://localhost/v2/keys/greeting", nil)
	recorder := httptest.NewRecorder()
	proxyHandler.ServeHTTP(recorder, request)

	if recorder.Code != 501 {
		t.Errorf("unexpected response code: %d", recorder.Code)
	}
}
