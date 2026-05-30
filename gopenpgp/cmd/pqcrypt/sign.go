package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
)

func runSign(args []string) {
	fs := flag.NewFlagSet("sign", flag.ExitOnError)
	key := fs.String("key", "", "Private key file (.key.asc) (required)")
	input := fs.String("input", "", "Input file to sign (reads stdin if not set)")
	output := fs.String("output", "", "Output signature file (default: <input>.sig)")
	fs.Parse(args)

	if *key == "" {
		fmt.Fprintln(os.Stderr, "Error: --key is required")
		fmt.Fprintln(os.Stderr, "\nUsage: pqcrypt sign --key alice.key.asc --input document.pdf")
		os.Exit(1)
	}

	var data []byte
	var err error
	if *input != "" {
		data, err = os.ReadFile(*input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			os.Exit(1)
		}
	} else {
		data, err = readStdin()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
			os.Exit(1)
		}
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

	pgp := crypto.PGP()

	signer, err := pgp.Sign().
		SigningKey(privKey).
		Detached().
		New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating signer: %v\n", err)
		os.Exit(1)
	}

	signature, err := signer.Sign(data, crypto.Armor)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error signing: %v\n", err)
		os.Exit(1)
	}

	outFile := *output
	if outFile == "" && *input != "" {
		outFile = *input + ".sig"
	}

	if outFile != "" {
		if err := os.WriteFile(outFile, signature, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing signature: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Signed: %s -> %s\n", *input, outFile)
		fmt.Printf("  Key: %s\n", privKey.GetFingerprint())
	} else {
		os.Stdout.Write(signature)
	}
}
