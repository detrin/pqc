package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
	"github.com/ProtonMail/gopenpgp/v3/profile"
)

func runKeygen(args []string) {
	fs := flag.NewFlagSet("keygen", flag.ExitOnError)
	name := fs.String("name", "", "User name (required)")
	email := fs.String("email", "", "User email (required)")
	pqc := fs.Bool("pqc", false, "Generate post-quantum key (ML-DSA + ML-KEM)")
	output := fs.String("output", "", "Output file prefix (default: <email>)")
	fs.Parse(args)

	if *name == "" || *email == "" {
		fmt.Fprintln(os.Stderr, "Error: --name and --email are required")
		fmt.Fprintln(os.Stderr, "\nUsage: pqcrypt keygen --name \"Alice\" --email alice@example.com [--pqc]")
		os.Exit(1)
	}

	var pgp *crypto.PGPHandle
	if *pqc {
		pgp = crypto.PGPWithProfile(profile.PQC())
		fmt.Println("Generating post-quantum key pair (ML-DSA-65+Ed25519 / ML-KEM-768+X25519)...")
	} else {
		pgp = crypto.PGPWithProfile(profile.RFC9580())
		fmt.Println("Generating classical key pair (Ed25519 / X25519)...")
	}

	key, err := pgp.KeyGeneration().
		AddUserId(*name, *email).
		New().
		GenerateKey()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating key: %v\n", err)
		os.Exit(1)
	}

	prefix := *output
	if prefix == "" {
		prefix = *email
	}

	privArmored, err := key.Armor()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error armoring private key: %v\n", err)
		os.Exit(1)
	}

	privFile := prefix + ".key.asc"
	if err := os.WriteFile(privFile, []byte(privArmored), 0600); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing private key: %v\n", err)
		os.Exit(1)
	}

	pubKey, err := key.ToPublic()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error extracting public key: %v\n", err)
		os.Exit(1)
	}

	pubArmored, err := pubKey.Armor()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error armoring public key: %v\n", err)
		os.Exit(1)
	}

	pubFile := prefix + ".pub.asc"
	if err := os.WriteFile(pubFile, []byte(pubArmored), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing public key: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nKey generated successfully!\n")
	fmt.Printf("  Fingerprint: %s\n", key.GetFingerprint())
	fmt.Printf("  Private key: %s\n", privFile)
	fmt.Printf("  Public key:  %s\n", pubFile)
	if *pqc {
		fmt.Println("\n  [!] This key uses post-quantum algorithms (draft-ietf-openpgp-pqc).")
		fmt.Println("      It is only interoperable with PQC-capable implementations.")
	}
}
