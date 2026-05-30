package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
)

func runDecrypt(args []string) {
	fs := flag.NewFlagSet("decrypt", flag.ExitOnError)
	key := fs.String("key", "", "Private key file (.key.asc)")
	input := fs.String("input", "", "Input file to decrypt (reads stdin if not set)")
	output := fs.String("output", "", "Output file (default: strips .pgp extension or stdout)")
	password := fs.String("password", "", "Decrypt with password instead of private key")
	fs.Parse(args)

	if *key == "" && *password == "" {
		fmt.Fprintln(os.Stderr, "Error: --key or --password is required")
		fmt.Fprintln(os.Stderr, "\nUsage: pqcrypt decrypt --key alice.key.asc --input message.txt.pgp")
		os.Exit(1)
	}

	var ciphertext []byte
	var err error
	if *input != "" {
		ciphertext, err = os.ReadFile(*input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			os.Exit(1)
		}
	} else {
		ciphertext, err = readStdin()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
			os.Exit(1)
		}
	}

	pgp := crypto.PGP()

	var decBuilder *crypto.DecryptionHandleBuilder
	if *password != "" {
		decBuilder = pgp.Decryption().Password([]byte(*password))
	} else {
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
		decBuilder = pgp.Decryption().DecryptionKey(privKey)
	}

	decHandle, err := decBuilder.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating decryption handle: %v\n", err)
		os.Exit(1)
	}

	decrypted, err := decHandle.Decrypt(ciphertext, crypto.Armor)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error decrypting: %v\n", err)
		os.Exit(1)
	}

	outFile := *output
	if outFile == "" && *input != "" {
		outFile = strings.TrimSuffix(*input, ".pgp")
		if outFile == *input {
			outFile = *input + ".dec"
		}
	}

	if outFile != "" {
		if err := os.WriteFile(outFile, decrypted.Bytes(), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing output: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Decrypted: %s -> %s\n", *input, outFile)
	} else {
		os.Stdout.Write(decrypted.Bytes())
	}
}
