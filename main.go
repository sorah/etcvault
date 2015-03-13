package main

import (
	"bufio"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/sorah/etcvault/engine"
	"github.com/sorah/etcvault/keys"
	"io"
	"net/url"
	"os"
	"strings"
	"time"
)

func main() {
	app := cli.NewApp()
	app.Name = "etcvault"
	app.Usage = "proxy for etcd, adding transparent encryption"

	app.Commands = []cli.Command{
		{
			Name:   "start",
			Usage:  "start etcvault proxy",
			Action: actionStart,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "keychain",
					Usage: "Path to directory for keys",
				},
				cli.StringFlag{
					Name:  "listen",
					Value: "http://localhost:2381",
					Usage: "URL to listen. Specify https as scheme to listen HTTPS.",
				},
				cli.StringFlag{
					Name:  "advertise-url",
					Value: "http://localhost:2381",
					Usage: "Client URL to advertise. Usually specify etcvault's URL",
				},

				cli.StringFlag{
					Name:  "discovery-srv",
					Usage: "domain to fetch SRV records for backend etcd",
				},
				cli.StringFlag{
					Name:  "initial-backends",
					Usage: "backend urls to fetch backend etcd members, separeted by comma",
				},
				cli.StringFlag{
					Name:  "client-ca-file",
					Usage: "TLS CA file to verify certificate of etcd client ports (https://...:2379/)",
				},
				cli.StringFlag{
					Name:  "client-cert-file",
					Usage: "TLS certficate file to send when communicating with etcd client ports (https://...:2379/)",
				},
				cli.StringFlag{
					Name:  "client-key-file",
					Usage: "key for -client-cert-file",
				},
				cli.StringFlag{
					Name:  "peer-ca-file",
					Usage: "TLS CA file to verify certificate of etcd peer ports (https://...:2380/)",
				},
				cli.StringFlag{
					Name:  "peer-cert-file",
					Usage: "TLS certficate file to send when communicating with etcd peer ports (https://...:2380/)",
				},
				cli.StringFlag{
					Name:  "peer-key-file",
					Usage: "key for -peer-cert-file",
				},
				cli.StringFlag{
					Name:  "listen-ca-file",
					Usage: "When listening HTTPS and this is present, etcvault will validate its client with using this CA certificate. If not present, -client-ca-file will be used.",
				},
				cli.StringFlag{
					Name:  "listen-cert-file",
					Usage: "When listening HTTPS and this is present, etcvault will use this certificate to listen. If not present, -client-cert-file will be used.",
				},
				cli.StringFlag{
					Name:  "listen-key-file",
					Usage: "key for -listen-cert-file",
				},
				cli.IntFlag{
					Name:  "discovery-interval",
					Value: 120,
					Usage: "Interval (in second) to refresh backends with specified discovery method",
				},
				cli.BoolFlag{
					Name:  "readonly",
					Usage: "if set, etcvault will reject non GET requests",
				},
			},
		},
		{
			Name:   "keygen",
			Usage:  "Generate new private key with specified name",
			Action: actionKeygen,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "save",
					Usage: "Save generated key into specfied directory (keychain)",
				},
				cli.IntFlag{
					Name:  "bits",
					Value: 2048,
					Usage: "RSA key bit length to generate",
				},
			},
		},
		{
			Name:   "transform",
			Usage:  "transform ETCVAULT* strings (from argument or stdin) to appropriate strings",
			Action: actionTransform,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "keychain",
					Usage: "Path to directory for keys",
				},
				cli.BoolFlag{
					Name:  "stdin",
					Usage: "Read from stdin",
				},
			},
		},
	}

	app.Run(os.Args)
}

func actionKeygen(ctx *cli.Context) {
	if len(ctx.Args()) < 1 {
		fmt.Fprintln(os.Stderr, "specify key name")
		os.Exit(1)
	}

	name := ctx.Args()[0]
	bits := ctx.Int("bits")

	key, err := keys.GenerateKey(name, bits)
	if err != nil {
		panic(err)
	}

	saveDir := ctx.String("save")

	if saveDir == "" {
		fmt.Printf("%s", key.PrivatePem())
	} else {
		keychain := keys.NewKeychain(saveDir)
		keychain.Save(key)
	}
}

func actionTransform(ctx *cli.Context) {
	keychainDir := ctx.String("keychain")
	if keychainDir == "" {
		fmt.Fprintln(os.Stderr, "Specify -keychain option")
		os.Exit(1)
	}

	keychain := keys.NewKeychain(keychainDir)
	engine := engine.NewEngine(keychain)

	if ctx.Bool("stdin") {
		reader := bufio.NewReader(os.Stderr)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					break
				} else {
					panic(err)
				}
			}

			origStr := strings.TrimRight(line, "\n")

			str, err := engine.Transform(origStr)
			if err == nil {
				fmt.Println(str)
			} else {
				fmt.Println(origStr)
				fmt.Fprintf(os.Stderr, "ERR: %s\n", err.Error())
			}
		}
	} else {
		for _, origStr := range ctx.Args() {
			str, err := engine.Transform(origStr)
			if err == nil {
				fmt.Println(str)
			} else {
				fmt.Println(origStr)
				fmt.Fprintf(os.Stderr, "ERR: %s", err.Error())
			}
		}
	}
}

func actionStart(ctx *cli.Context) {
	keychainDir := ctx.String("keychain")
	if keychainDir == "" {
		fmt.Fprintln(os.Stderr, "Specify -keychain option")
		os.Exit(1)
	}

	discoverySrvDomain := ctx.String("discovery-srv")
	initialBackendUrlStrings := ctx.String("initial-backends")
	if discoverySrvDomain == "" && initialBackendUrlStrings == "" {
		fmt.Fprintln(os.Stderr, "Specify -discovery-srv or -initial-backends option")
		os.Exit(1)
	}
	if discoverySrvDomain != "" && initialBackendUrlStrings != "" {
		fmt.Fprintln(os.Stderr, "Only specifying only either -discovery-srv or -initial-backends is accepted.")
		os.Exit(1)
	}

	clientCaFilePath := ctx.String("client-ca-file")
	clientCertFilePath := ctx.String("client-cert-file")
	clientKeyFilePath := ctx.String("client-key-file")
	if (clientCertFilePath != "" || clientKeyFilePath != "") && !(clientCertFilePath != "" && clientKeyFilePath != "") {
		fmt.Fprintln(os.Stderr, "provide both -client-cert-file and -client-key-file")
		os.Exit(1)
	}

	peerCaFilePath := ctx.String("peer-ca-file")
	peerCertFilePath := ctx.String("peer-cert-file")
	peerKeyFilePath := ctx.String("peer-key-file")
	if (peerCertFilePath != "" || peerKeyFilePath != "") && !(peerCertFilePath != "" && peerKeyFilePath != "") {
		fmt.Fprintln(os.Stderr, "provide both -peer-cert-file and -peer-key-file")
		os.Exit(1)
	}

	listenCaFilePath := ctx.String("listen-ca-file")
	listenCertFilePath := ctx.String("listen-cert-file")
	listenKeyFilePath := ctx.String("listen-key-file")
	if (listenCertFilePath != "" || listenKeyFilePath != "") && !(listenCertFilePath != "" && listenKeyFilePath != "") {
		fmt.Fprintln(os.Stderr, "provide both -listen-cert-file and -listen-key-file")
		os.Exit(1)
	}

	discoveryInterval := ctx.Int("discovery-interval")

	readonly := ctx.Bool("readonly")

	listenUrl, err := url.Parse(ctx.String("listen"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't parse -listen as URL: %s\n", err.Error())
		os.Exit(1)
	}
	if listenUrl.Path != "" && listenUrl.Path != "/ " {
		fmt.Fprintf(os.Stderr, "-listen URL shouldn't include path: %s\n", listenUrl.Path)
		os.Exit(1)
	}
	if !(clientCertFilePath != "" && clientKeyFilePath != "") && listenUrl.Scheme == "https" {
		fmt.Fprintln(os.Stderr, "provide both -cert-file and -key-file when listen https")
		os.Exit(1)
	}

	advertiseUrl := ctx.String("advertise-url")

	starter := &ProxyStarter{
		Listen:                   listenUrl,
		keychainDir:              keychainDir,
		DiscoverySrvDomain:       discoverySrvDomain,
		initialBackendUrlStrings: initialBackendUrlStrings,
		clientCaFilePath:         clientCaFilePath,
		clientCertFilePath:       clientCertFilePath,
		clientKeyFilePath:        clientKeyFilePath,
		peerCaFilePath:           peerCaFilePath,
		peerCertFilePath:         peerCertFilePath,
		peerKeyFilePath:          peerKeyFilePath,
		listenCaFilePath:         listenCaFilePath,
		listenCertFilePath:       listenCertFilePath,
		listenKeyFilePath:        listenKeyFilePath,
		discoveryInterval:        time.Duration(discoveryInterval) * time.Second,
		readonly:                 readonly,
		AdvertiseUrl:             advertiseUrl,
	}

	starter.Start()
}
