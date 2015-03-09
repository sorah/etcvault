package main

import (
	"bufio"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/sorah/etcvault/engine"
	"github.com/sorah/etcvault/keys"
	"io"
	"os"
	"strings"
)

func main() {
	app := cli.NewApp()
	app.Name = "etcvault"
	app.Usage = "proxy for etcd, adding transparent encryption"

	app.Commands = []cli.Command{
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
