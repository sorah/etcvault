package proxy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
)

// for testing...
var lookupSRV = net.LookupSRV

type etcdMember struct {
	ClientURLs []string
	PeerURLs   []string
	Name       string
	Id         string
}

type etcdMembers struct {
	Members []etcdMember
}

func DiscoverBackendsFromDns(transport *http.Transport, domain string) ([]*Backend, error) {
	_, records, errA := lookupSRV("etcd-server", "tcp", domain)
	if errA != nil {
		log.Printf("error when looking up _etcd-server._tcp.%s: %s", domain, errA.Error())
	}

	_, ssl_records, errB := lookupSRV("etcd-server-ssl", "tcp", domain)
	if errB != nil {
		log.Printf("error when looking up _etcd-server-ssl._tcp.%s: %s", domain, errB.Error())
	}

	if errA != nil && errB != nil {
		return nil, errA
	}

	urls := make([]*url.URL, 0, len(records)+len(ssl_records))

	makeUrl := func(srv *net.SRV, scheme string) *url.URL {
		var target string
		if srv.Target[len(srv.Target)-1] == '.' {
			target = srv.Target[0 : len(srv.Target)-1]
		} else {
			target = srv.Target
		}

		hostPort := net.JoinHostPort(target, fmt.Sprintf("%d", srv.Port))

		u := &url.URL{
			Scheme: scheme,
			Host:   hostPort,
		}
		return u
	}

	for _, srv := range ssl_records {
		urls = append(urls, makeUrl(srv, "https"))
	}

	for _, srv := range records {
		urls = append(urls, makeUrl(srv, "http"))
	}

	return DiscoverBackendsFromEtcdPeer(transport, urls), nil
}

func DiscoverBackendsFromEtcdPeer(transport *http.Transport, urls []*url.URL) []*Backend {
	return fetchBackendsFromEtcd(transport, urls, "/members")
}

func DiscoverBackendsFromEtcd(transport *http.Transport, urls []*url.URL) []*Backend {
	return fetchBackendsFromEtcd(transport, urls, "/v2/members")
}

func fetchBackendsFromEtcd(transport *http.Transport, urls []*url.URL, path string) []*Backend {
	client := &http.Client{Transport: transport}

	for _, origUrl := range urls {
		u := new(url.URL)
		*u = *origUrl

		u.Path = path

		resp, err := client.Get(u.String())
		if err != nil {
			log.Printf("error when retrieving %s: %s", u.String(), err.Error())
			continue
		}

		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			continue
		}
		err = resp.Body.Close()
		if err != nil {
			panic(err)
		}

		members := &etcdMembers{}
		err = json.Unmarshal(respBody, members)
		if err != nil {
			continue
		}

		backends := make([]*Backend, 0, len(members.Members))

		for _, member := range members.Members {
			if len(member.ClientURLs) < 1 {
				continue
			}
			clientUrl, err := url.Parse(member.ClientURLs[0])
			if err != nil {
				continue
			}
			backend := NewBackend(clientUrl)

			backends = append(backends, backend)
		}

		return backends
	}

	return []*Backend{}
}
