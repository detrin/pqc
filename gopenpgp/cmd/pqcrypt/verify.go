package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
)

func runVerify(args []string) {
	fs := flag.NewFlagSet("verify", flag.ExitOnError)
	pubkey := fs.String("pubkey", "", "Public key file (.pub.asc) (required)")
	input := fs.String("input", "", "Input file that was signed (required)")
	signature := fs.String("signature", "", "Signature file (default: <input>.sig)")
	fs.Parse(args)

	if *pubkey == "" || *input == "" {
		fmt.Fprintln(os.Stderr, "Error: --pubkey and --input are required")
		fmt.Fprintln(os.Stderr, "\nUsage: pqcrypt verify --pubkey alice.pub.asc --input document.pdf")
		os.Exit(1)
	}

	sigFile := *signature
	if sigFile == "" {
		sigFile = *input + ".sig"
	}

	data, err := os.ReadFile(*input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	sigData, err := os.ReadFile(sigFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading signature: %v\n", err)
		os.Exit(1)
	}

	keyData, err := os.ReadFile(*pubkey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading public key: %v\n", err)
		os.Exit(1)
	}

	pubKey, err := crypto.NewKeyFromArmored(string(keyData))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing public key: %v\n", err)
		os.Exit(1)
	}

	pgp := crypto.PGP()

	verifier, err := pgp.Verify().
		VerificationKey(pubKey).
		New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating verifier: %v\n", err)
		os.Exit(1)
	}

	result, err := verifier.VerifyDetached(data, sigData, crypto.Armor)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error verifying: %v\n", err)
		os.Exit(1)
	}

	if sigErr := result.SignatureError(); sigErr != nil {
		fmt.Fprintf(os.Stderr, "INVALID signature: %v\n", sigErr)
		os.Exit(1)
	}

	fmt.Println("VALID signature.")
	fmt.Printf("  File:   %s\n", *input)
	fmt.Printf("  Signer: %s\n", pubKey.GetFingerprint())
}
