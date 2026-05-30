package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
)

func runExport(args []string) {
	fs := flag.NewFlagSet("export", flag.ExitOnError)
	key := fs.String("key", "", "Private key file (.key.asc) (required)")
	output := fs.String("output", "", "Output public key file (default: stdout)")
	fs.Parse(args)

	if *key == "" {
		fmt.Fprintln(os.Stderr, "Error: --key is required")
		fmt.Fprintln(os.Stderr, "\nUsage: pqcrypt export --key alice.key.asc")
		os.Exit(1)
	}

	keyData, err := os.ReadFile(*key)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading key: %v\n", err)
		os.Exit(1)
	}

	privKey, err := crypto.NewKeyFromArmored(string(keyData))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing key: %v\n", err)
		os.Exit(1)
	}

	pubKey, err := privKey.ToPublic()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error extracting public key: %v\n", err)
		os.Exit(1)
	}

	pubArmored, err := pubKey.Armor()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error armoring public key: %v\n", err)
		os.Exit(1)
	}

	if *output != "" {
		if err := os.WriteFile(*output, []byte(pubArmored), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing public key: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Exported public key: %s\n", *output)
	} else {
		fmt.Print(pubArmored)
	}
}
