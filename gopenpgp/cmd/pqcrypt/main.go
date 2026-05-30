package main

import (
	"fmt"
	"os"
)

const usage = `pqcrypt - Post-Quantum PGP Tool

Usage:
  pqcrypt <command> [options]

Commands:
  keygen      Generate a new PGP key pair (classical or post-quantum)
  encrypt     Encrypt a file or message for a recipient
  decrypt     Decrypt a file or message
  sign        Sign a file or message
  verify      Verify a signature
  export      Export public key from a private key
  inspect     Show key information

Options:
  --help, -h  Show this help message

Examples:
  pqcrypt keygen --name "Alice" --email alice@example.com --pqc
  pqcrypt encrypt --recipient alice.pub.asc --input secret.txt
  pqcrypt decrypt --key alice.key.asc --input secret.txt.pgp
  pqcrypt sign --key alice.key.asc --input document.pdf
  pqcrypt verify --pubkey alice.pub.asc --input document.pdf --signature document.pdf.sig
  pqcrypt export --key alice.key.asc
  pqcrypt inspect --key alice.pub.asc
`

func main() {
	if len(os.Args) < 2 {
		fmt.Print(usage)
		os.Exit(0)
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "keygen":
		runKeygen(args)
	case "encrypt":
		runEncrypt(args)
	case "decrypt":
		runDecrypt(args)
	case "sign":
		runSign(args)
	case "verify":
		runVerify(args)
	case "export":
		runExport(args)
	case "inspect":
		runInspect(args)
	case "--help", "-h", "help":
		fmt.Print(usage)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", cmd)
		fmt.Print(usage)
		os.Exit(1)
	}
}
