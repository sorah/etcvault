package proxy

import (
	"fmt"
	"net/url"
	"testing"
	"time"
)

func generateBackendsForTest(count int) []*Backend {
	backends := make([]*Backend, 0, count)

	for i := 0; i < count; i++ {
		u, _ := url.Parse(fmt.Sprintf("http://backend-%d", i))
		backends = append(backends, NewBackend(u))
	}

	return backends
}

func TestBackends(t *testing.T) {
	router := NewRouter(time.Second*60, func() ([]*Backend, error) {
		return generateBackendsForTest(3), nil
	})
	router.Update()

	backends := router.Backends()

	if len(backends) != 3 {
		t.Errorf("Unexpected backends length %d", len(backends))
		return
	}

	if backends[0].Url.Host != "backend-0" {
		t.Errorf("Unexpected backends[0] url %s", backends[0].Url.Host)
	}
	if backends[1].Url.Host != "backend-1" {
		t.Errorf("Unexpected backends[1] url %s", backends[1].Url.Host)
	}
	if backends[2].Url.Host != "backend-2" {
		t.Errorf("Unexpected backends[2] url %s", backends[2].Url.Host)
	}
}

func TestAvailableBackends(t *testing.T) {
	backendsSource := generateBackendsForTest(3)
	router := NewRouter(time.Second*60, func() ([]*Backend, error) {
		return backendsSource, nil
	})
	router.Update()

	backendsSource[0].Fail()

	backends := router.AvailableBackends()

	if len(backends) != 2 {
		t.Errorf("Unexpected backends length %d", len(backends))
		return
	}

	if backends[0].Url.Host != "backend-1" {
		t.Errorf("Unexpected backends[0] url %s", backends[0].Url.Host)
	}
	if backends[1].Url.Host != "backend-2" {
		t.Errorf("Unexpected backends[1] url %s", backends[1].Url.Host)
	}
}

func TestFailedBackends(t *testing.T) {
	backendsSource := generateBackendsForTest(3)
	router := NewRouter(time.Second*60, func() ([]*Backend, error) {
		return backendsSource, nil
	})
	router.Update()

	backendsSource[1].Fail()
	backendsSource[2].Fail()

	backends := router.FailedBackends()

	if len(backends) != 2 {
		t.Errorf("Unexpected backends length %d", len(backends))
		return
	}

	if backends[0].Url.Host != "backend-1" {
		t.Errorf("Unexpected backends[0] url %s", backends[0].Url.Host)
	}
	if backends[1].Url.Host != "backend-2" {
		t.Errorf("Unexpected backends[1] url %s", backends[1].Url.Host)
	}
}

func TestShuffledAvailableBackends(t *testing.T) {
	router := NewRouter(time.Second*60, func() ([]*Backend, error) {
		return generateBackendsForTest(3), nil
	})
	router.Update()

	backends := router.ShuffledAvailableBackends()

	if len(backends) != 3 {
		t.Errorf("Unexpected backends length %d", len(backends))
		return
	}

	hosts := make(map[string]bool)
	hosts[backends[0].Url.Host] = true
	hosts[backends[1].Url.Host] = true
	hosts[backends[2].Url.Host] = true

	if exist, ok := hosts["backend-0"]; !(ok && exist) {
		t.Errorf("backend-0 not ok: %#v", backends)
	}
	if exist, ok := hosts["backend-1"]; !(ok && exist) {
		t.Errorf("backend-1 not ok: %#v", backends)
	}
	if exist, ok := hosts["backend-2"]; !(ok && exist) {
		t.Errorf("backend-2 not ok: %#v", backends)
	}
}

func TestUpdate(t *testing.T) {
	i := 0
	router := NewRouter(time.Second*60, func() ([]*Backend, error) {
		i++
		return generateBackendsForTest(2 + i), nil
	})

	router.Update()
	backends := router.Backends()

	if len(backends) != 3 {
		t.Errorf("Unexpected backends length %d", len(backends))
		return
	}

	if backends[0].Url.Host != "backend-0" {
		t.Errorf("Unexpected backends[0] url %s", backends[0].Url.Host)
	}
	if backends[1].Url.Host != "backend-1" {
		t.Errorf("Unexpected backends[1] url %s", backends[1].Url.Host)
	}
	if backends[2].Url.Host != "backend-2" {
		t.Errorf("Unexpected backends[2] url %s", backends[2].Url.Host)
	}

	router.Update()
	backends = router.Backends()

	if len(backends) != 4 {
		t.Errorf("Unexpected backends length %d", len(backends))
		return
	}

	if backends[0].Url.Host != "backend-0" {
		t.Errorf("Unexpected backends[0] url %s", backends[0].Url.Host)
	}
	if backends[1].Url.Host != "backend-1" {
		t.Errorf("Unexpected backends[1] url %s", backends[1].Url.Host)
	}
	if backends[2].Url.Host != "backend-2" {
		t.Errorf("Unexpected backends[2] url %s", backends[2].Url.Host)
	}
	if backends[3].Url.Host != "backend-3" {
		t.Errorf("Unexpected backends[3] url %s", backends[3].Url.Host)
	}
}

func TestUpdateFail(t *testing.T) {
	i := -1
	router := NewRouter(time.Second*60, func() (backends []*Backend, err error) {
		i++
		if i == 1 {
			backends = nil
			err = fmt.Errorf("hehe")
			return
		}
		backends = generateBackendsForTest(3 + i)
		return
	})

	router.Update()
	router.Update()
	backends := router.Backends()

	if len(backends) != 3 {
		t.Errorf("Unexpected backends length %d", len(backends))
		return
	}

	if backends[0].Url.Host != "backend-0" {
		t.Errorf("Unexpected backends[0] url %s", backends[0].Url.Host)
	}
	if backends[1].Url.Host != "backend-1" {
		t.Errorf("Unexpected backends[1] url %s", backends[1].Url.Host)
	}
	if backends[2].Url.Host != "backend-2" {
		t.Errorf("Unexpected backends[2] url %s", backends[2].Url.Host)
	}

	router.Update()
	backends = router.Backends()

	if len(backends) != 5 {
		t.Errorf("Unexpected backends length %d", len(backends))
		return
	}

	if backends[0].Url.Host != "backend-0" {
		t.Errorf("Unexpected backends[0] url %s", backends[0].Url.Host)
	}
	if backends[1].Url.Host != "backend-1" {
		t.Errorf("Unexpected backends[1] url %s", backends[1].Url.Host)
	}
	if backends[2].Url.Host != "backend-2" {
		t.Errorf("Unexpected backends[2] url %s", backends[2].Url.Host)
	}
	if backends[3].Url.Host != "backend-3" {
		t.Errorf("Unexpected backends[3] url %s", backends[3].Url.Host)
	}
	if backends[4].Url.Host != "backend-4" {
		t.Errorf("Unexpected backends[4] url %s", backends[4].Url.Host)
	}
}
