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
	_, records, err := lookupSRV("etcd-server", "tcp", domain)

	if err != nil {
		return nil, err
	}

	urls := make([]*url.URL, 0, len(records))
	for _, srv := range records {
		var target string
		if srv.Target[len(srv.Target)-1] == '.' {
			target = srv.Target[0 : len(srv.Target)-1]
		} else {
			target = srv.Target
		}

		hostPort := net.JoinHostPort(target, fmt.Sprintf("%d", srv.Port))

		u := &url.URL{
			Scheme: "http",
			Host:   hostPort,
		}
		urls = append(urls, u)
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
