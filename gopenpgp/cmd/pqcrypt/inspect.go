package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
)

func algoName(id uint8) string {
	names := map[uint8]string{
		1:  "RSA",
		2:  "RSA (Encrypt Only)",
		3:  "RSA (Sign Only)",
		16: "Elgamal",
		17: "DSA",
		18: "ECDH",
		19: "ECDSA",
		22: "EdDSA (Legacy)",
		25: "X25519",
		27: "Ed25519",
		26: "X448",
		28: "Ed448",
		29: "ML-KEM-768+X25519",
		30: "ML-KEM-1024+X448",
		31: "ML-DSA-65+Ed25519",
		32: "ML-DSA-87+Ed448",
		33: "SLH-DSA-128s",
		34: "SLH-DSA-128f",
		35: "SLH-DSA-256s",
	}
	if name, ok := names[id]; ok {
		return name
	}
	return fmt.Sprintf("Unknown(%d)", id)
}

func runInspect(args []string) {
	fs := flag.NewFlagSet("inspect", flag.ExitOnError)
	keyFile := fs.String("key", "", "Key file to inspect (.key.asc or .pub.asc) (required)")
	fs.Parse(args)

	if *keyFile == "" {
		fmt.Fprintln(os.Stderr, "Error: --key is required")
		fmt.Fprintln(os.Stderr, "\nUsage: pqcrypt inspect --key alice.pub.asc")
		os.Exit(1)
	}

	keyData, err := os.ReadFile(*keyFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading key: %v\n", err)
		os.Exit(1)
	}

	key, err := crypto.NewKeyFromArmored(string(keyData))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing key: %v\n", err)
		os.Exit(1)
	}

	keyType := "Public"
	if key.IsPrivate() {
		keyType = "Private"
	}

	entity := key.GetEntity()
	var identities []string
	for id := range entity.Identities {
		identities = append(identities, id)
	}

	fmt.Printf("Key Information:\n")
	fmt.Printf("  Type:         %s Key\n", keyType)
	fmt.Printf("  Fingerprint:  %s\n", key.GetFingerprint())
	fmt.Printf("  Identities:   %s\n", strings.Join(identities, ", "))
	fmt.Printf("  Created:      %s\n", entity.PrimaryKey.CreationTime.Format("2006-01-02 15:04:05"))

	algo := entity.PrimaryKey.PubKeyAlgo
	fmt.Printf("  Algorithm:    %s (ID: %d)\n", algoName(uint8(algo)), algo)

	subkeys := entity.Subkeys
	fmt.Printf("  Subkeys:      %d\n", len(subkeys))
	for i, sub := range subkeys {
		subAlgo := sub.PublicKey.PubKeyAlgo
		fmt.Printf("    [%d] %s (ID: %d) - fingerprint: %x\n",
			i+1, algoName(uint8(subAlgo)), subAlgo, sub.PublicKey.Fingerprint)
	}
}
