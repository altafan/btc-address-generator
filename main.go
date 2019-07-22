package main

import (
	"fmt"
	"os"
)

func exit(err string) {
	fmt.Println(err)
	os.Exit(1)
}

func main() {
	xpub, ok := os.LookupEnv("XPUB")
	if !ok {
		exit("Invalid xpub: must not be empty")
	}

	network, ok := os.LookupEnv("NETWORK")
	if !ok {
		exit("Invalid network: must not be empty")
	}

	derivationPath, ok := os.LookupEnv("DERIVATION_PATH")
	if !ok {
		exit("Invalid account index: must not be empty")
	}

	l, err := NewLambda(network, xpub, derivationPath)
	if err != nil {
		exit(err.Error())
	}

	l.Start()
}
