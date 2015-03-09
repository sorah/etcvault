package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/sorah/etcvault/keys"
	"os"
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
