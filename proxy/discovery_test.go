package proxy

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func membersMock(count int, peer bool) *httptest.Server {
	type memberT struct {
		ClientURLs []string
		PeerURLs   []string
		Name       string
		Id         string
	}

	members := make([]memberT, 0, count)
	for i := 0; i < count; i++ {
		member := memberT{
			ClientURLs: []string{fmt.Sprintf("http://member-%d:2379", i)},
			PeerURLs:   []string{fmt.Sprintf("http://member-%d:2380", i)},
			Name:       fmt.Sprintf("member-%d", i),
			Id:         fmt.Sprintf("%x", i),
		}
		members = append(members, member)
	}

	membersJson, err := json.Marshal(struct {
		Members []memberT
	}{
		Members: members,
	})
	if err != nil {
		panic(err)
	}

	var path, host string
	if peer {
		host = "node:2380"
		path = "/members"
	} else {
		host = "node:2379"
		path = "/v2/members"
	}

	server := httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, request *http.Request) {
		if request.Host == host && request.URL.Path == path && request.Method == "GET" {
			resp.Header().Add("Content-Type", "application/json")
			resp.WriteHeader(200)
			_, _ = resp.Write(membersJson)
		} else {
			http.Error(resp, "not found", 404)
		}
	}))

	return server
}

// ----

func TestDiscoverBackendsFromEtcd(t *testing.T) {
	testServer := membersMock(3, false)
	defer testServer.Close()

	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			u, _ := url.Parse(testServer.URL)
			return u, nil
		},
	}

	u, err := url.Parse("http://node:2379")
	if err != nil {
		panic(err)
	}

	backends := DiscoverBackendsFromEtcd(transport, []*url.URL{u})

	if len(backends) != 3 {
		t.Errorf("unexpected backends size %d", len(backends))
		return
	}

	if backends[0].Url.String() != "http://member-0:2379" {
		t.Errorf("unexpected backends[0] url %s", backends[0].Url.String())
	}
	if backends[1].Url.String() != "http://member-1:2379" {
		t.Errorf("unexpected backends[1] url %s", backends[1].Url.String())
	}
	if backends[2].Url.String() != "http://member-2:2379" {
		t.Errorf("unexpected backends[2] url %s", backends[2].Url.String())
	}

	if u.String() != "http://node:2379" {
		t.Errorf("url changed %s", u.String())
	}
}

func TestDiscoverBackendsFromEtcdPeer(t *testing.T) {
	testServer := membersMock(3, true)
	defer testServer.Close()

	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			u, _ := url.Parse(testServer.URL)
			return u, nil
		},
	}

	u, err := url.Parse("http://node:2380")
	if err != nil {
		panic(err)
	}

	backends := DiscoverBackendsFromEtcdPeer(transport, []*url.URL{u})

	if len(backends) != 3 {
		t.Errorf("unexpected backends size %d", len(backends))
		return
	}

	if backends[0].Url.String() != "http://member-0:2379" {
		t.Errorf("unexpected backends[0] url %s", backends[0].Url.String())
	}
	if backends[1].Url.String() != "http://member-1:2379" {
		t.Errorf("unexpected backends[1] url %s", backends[1].Url.String())
	}
	if backends[2].Url.String() != "http://member-2:2379" {
		t.Errorf("unexpected backends[2] url %s", backends[2].Url.String())
	}

	if u.String() != "http://node:2380" {
		t.Errorf("url changed %s", u.String())
	}
}

func TestDiscoverBackendsFromDns(t *testing.T) {
	testServer := membersMock(3, true)
	defer testServer.Close()

	lookupSRV = func(service, proto, name string) (string, []*net.SRV, error) {
		if service == "etcd-server" && proto == "tcp" && name == "example.org" {
			return "", []*net.SRV{
				{
					Target:   "node.",
					Port:     2380,
					Priority: 0,
					Weight:   0,
				},
			}, nil
		}
		return "", []*net.SRV{}, &net.DNSError{Err: "no such host", Name: "", Server: "", IsTimeout: false}
	}
	defer func() { lookupSRV = net.LookupSRV }()

	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			u, _ := url.Parse(testServer.URL)
			return u, nil
		},
	}

	backends, err := DiscoverBackendsFromDns(transport, "example.org")

	if err != nil {
		t.Errorf("err %s", err.Error())
	}

	if len(backends) != 3 {
		t.Errorf("unexpected backends size %d", len(backends))
		return
	}

	if backends[0].Url.String() != "http://member-0:2379" {
		t.Errorf("unexpected backends[0] url %s", backends[0].Url.String())
	}
	if backends[1].Url.String() != "http://member-1:2379" {
		t.Errorf("unexpected backends[1] url %s", backends[1].Url.String())
	}
	if backends[2].Url.String() != "http://member-2:2379" {
		t.Errorf("unexpected backends[2] url %s", backends[2].Url.String())
	}
}