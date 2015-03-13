package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/sorah/etcvault/engine"
	"github.com/sorah/etcvault/keys"
	"github.com/sorah/etcvault/proxy"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func defaultHttpTransport() *http.Transport {
	return &http.Transport{
		// DefaultTransport
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	}
}

func caPool(caPath string) *x509.CertPool {
	pool := x509.NewCertPool()
	remainingPem, err := ioutil.ReadFile(caPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading CA file %s: %s", caPath, err)
		os.Exit(1)
	}

	for { // load while file ends
		var block *pem.Block
		block, remainingPem = pem.Decode(remainingPem)
		if block == nil {
			return pool
		}
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error while parsing CA PEM blocks: %s", err.Error())
			os.Exit(1)
		}
		pool.AddCert(cert)
	}
}

func parseTlsKeypair(certPath, keyPath string) *tls.Config {
	certBytes, err := ioutil.ReadFile(certPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading certificate %s: %s\n", certPath, err.Error())
		os.Exit(1)
	}
	keyBytes, err := ioutil.ReadFile(keyPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading certificate %s: %s\n", certPath, err.Error())
		os.Exit(1)
	}

	keypair, err := tls.X509KeyPair(certBytes, keyBytes)
	if err != nil {
		fmt.Printf("error loading keypair: %s", err.Error())
	}

	return &tls.Config{
		Certificates: []tls.Certificate{keypair},
		MinVersion:   tls.VersionTLS10,
	}
}

func tlsConfigurationForClientUse(config *tls.Config, caPath string) *tls.Config {
	if config == nil {
		return nil
	}

	config.RootCAs = caPool(caPath)
	return config
}

func tlsConfigurationForServerUse(config *tls.Config, caPath string) *tls.Config {
	if config == nil {
		return nil
	}

	if caPath != "" {
		config.ClientAuth = tls.RequireAndVerifyClientCert
		config.ClientCAs = caPool(caPath)
	} else {
		config.ClientAuth = tls.NoClientCert
	}
	return config
}

type ProxyStarter struct {
	// arguments
	Listen *url.URL

	keychainDir              string
	DiscoverySrvDomain       string
	initialBackendUrlStrings string

	clientCaFilePath   string
	clientCertFilePath string
	clientKeyFilePath  string

	peerCaFilePath   string
	peerCertFilePath string
	peerKeyFilePath  string

	readonly bool

	discoveryInterval time.Duration

	router *proxy.Router
}

func (starter *ProxyStarter) InitialBackendUrls() []*url.URL {
	urlStrings := strings.Split(starter.initialBackendUrlStrings, ",")
	urls := make([]*url.URL, len(urlStrings))

	for i, urlString := range urlStrings {
		u, err := url.Parse(urlString)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse url %s: %s\n", urlString, err.Error())
			os.Exit(1)
		}
		urls[i] = u
	}

	return urls
}

func (starter *ProxyStarter) Keychain() *keys.Keychain {
	return keys.NewKeychain(starter.keychainDir)
}

func (starter *ProxyStarter) Engine() *engine.Engine {
	return engine.NewEngine(starter.Keychain())
}

func (starter *ProxyStarter) PeerTlsConfig() *tls.Config {
	if starter.peerKeyFilePath != "" && starter.peerCertFilePath != "" {
		return parseTlsKeypair(starter.peerCertFilePath, starter.peerKeyFilePath)
	} else {
		return nil
	}
}

func (starter *ProxyStarter) ClientTlsConfig() *tls.Config {
	if starter.clientKeyFilePath != "" && starter.clientCertFilePath != "" {
		return parseTlsKeypair(starter.clientCertFilePath, starter.clientKeyFilePath)
	} else {
		return nil
	}
}

func (starter *ProxyStarter) ClientTlsConfigForClientUse() *tls.Config {
	return tlsConfigurationForClientUse(starter.ClientTlsConfig(), starter.clientCaFilePath)
}

func (starter *ProxyStarter) ClientTlsConfigForServerUse() *tls.Config {
	return tlsConfigurationForServerUse(starter.ClientTlsConfig(), starter.clientCaFilePath)
}

func (starter *ProxyStarter) PeerTlsConfigForClientUse() *tls.Config {
	return tlsConfigurationForClientUse(starter.PeerTlsConfig(), starter.peerCaFilePath)
}

func (starter *ProxyStarter) PeerHttpTransport() *http.Transport {
	transport := defaultHttpTransport()
	transport.TLSClientConfig = starter.PeerTlsConfigForClientUse()
	return transport
}

func (starter *ProxyStarter) ClientHttpTransport() *http.Transport {
	transport := defaultHttpTransport()
	transport.TLSClientConfig = starter.ClientTlsConfigForClientUse()
	return transport
}

func (starter *ProxyStarter) Listener() net.Listener {
	listener, err := net.Listen("tcp", starter.Listen.Host)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to listen %s: %s", starter.Listen.String(), err.Error())
		os.Exit(1)
	}

	if starter.Listen.Scheme == "https" {
		tlsConfig := starter.ClientTlsConfigForServerUse()
		listener = tls.NewListener(listener, tlsConfig)
	}

	return listener
}

func (starter *ProxyStarter) BackendUpdateFunc() proxy.BackendUpdateFunc {
	if starter.DiscoverySrvDomain != "" {
		transport := starter.PeerHttpTransport()
		return func() ([]*proxy.Backend, error) {
			return proxy.DiscoverBackendsFromDns(transport, starter.DiscoverySrvDomain)
		}
	} else {
		transport := starter.ClientHttpTransport()
		return func() ([]*proxy.Backend, error) {
			return proxy.DiscoverBackendsFromEtcd(transport, starter.InitialBackendUrls()), nil
		}
	}
}

func (starter *ProxyStarter) Router() *proxy.Router {
	if starter.router != nil {
		return starter.router
	}

	starter.router = proxy.NewRouter(starter.discoveryInterval, starter.BackendUpdateFunc())
	err := starter.router.StartUpdate()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error starting backend discovery: %s", err.Error())
	}

	return starter.router
}

func (starter *ProxyStarter) Proxy() http.Handler {
	if starter.readonly {
		return proxy.NewReadonlyProxy(starter.ClientHttpTransport(), starter.Router(), starter.Engine())
	} else {
		return proxy.NewProxy(starter.ClientHttpTransport(), starter.Router(), starter.Engine())
	}
}

func (starter *ProxyStarter) HttpServer() *http.Server {
	return &http.Server{
		Handler:     starter.Proxy(),
		ReadTimeout: 5 * time.Minute,
	}
}

func (starter *ProxyStarter) Start() {
	fmt.Printf("Serving at %s\n", starter.Listen.String())
	starter.HttpServer().Serve(starter.Listener())
}
