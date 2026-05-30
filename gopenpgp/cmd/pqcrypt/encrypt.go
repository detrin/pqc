package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
)

func runEncrypt(args []string) {
	fs := flag.NewFlagSet("encrypt", flag.ExitOnError)
	recipient := fs.String("recipient", "", "Recipient public key file (.pub.asc)")
	input := fs.String("input", "", "Input file to encrypt (reads stdin if not set)")
	output := fs.String("output", "", "Output file (default: <input>.pgp or stdout)")
	password := fs.String("password", "", "Encrypt with password instead of public key")
	fs.Parse(args)

	if *recipient == "" && *password == "" {
		fmt.Fprintln(os.Stderr, "Error: --recipient or --password is required")
		fmt.Fprintln(os.Stderr, "\nUsage: pqcrypt encrypt --recipient bob.pub.asc --input message.txt")
		fmt.Fprintln(os.Stderr, "       pqcrypt encrypt --password --input secret.txt")
		os.Exit(1)
	}

	var plaintext []byte
	var err error
	if *input != "" {
		plaintext, err = os.ReadFile(*input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			os.Exit(1)
		}
	} else {
		plaintext, err = readStdin()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
			os.Exit(1)
		}
	}

	pgp := crypto.PGP()

	var encBuilder *crypto.EncryptionHandleBuilder
	if *password != "" {
		encBuilder = pgp.Encryption().Password([]byte(*password))
	} else {
		keyData, err := os.ReadFile(*recipient)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading recipient key: %v\n", err)
			os.Exit(1)
		}

		pubKey, err := crypto.NewKeyFromArmored(string(keyData))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing recipient key: %v\n", err)
			os.Exit(1)
		}
		encBuilder = pgp.Encryption().Recipient(pubKey)
	}

	encHandle, err := encBuilder.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating encryption handle: %v\n", err)
		os.Exit(1)
	}

	ciphertext, err := encHandle.Encrypt(plaintext)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error encrypting: %v\n", err)
		os.Exit(1)
	}

	armored, err := ciphertext.ArmorBytes()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error armoring ciphertext: %v\n", err)
		os.Exit(1)
	}

	outFile := *output
	if outFile == "" && *input != "" {
		outFile = *input + ".pgp"
	}

	if outFile != "" {
		if err := os.WriteFile(outFile, armored, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing output: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Encrypted: %s -> %s\n", *input, outFile)
	} else {
		os.Stdout.Write(armored)
	}
}
