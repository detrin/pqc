package main

import (
	"fmt"
	"log"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
	"github.com/ProtonMail/gopenpgp/v3/profile"
)

func main() {
	fmt.Println("=== GopenPGP Post-Quantum Cryptography Examples ===")
	fmt.Println()

	classicalExample()
	pqcExample()
	signAndVerifyExample()
	detachedSignatureExample()
	passwordEncryptionExample()
}

func classicalExample() {
	fmt.Println("--- 1. Classical Key Generation & Encryption (Ed25519/X25519) ---")

	pgp := crypto.PGP()

	key, err := pgp.KeyGeneration().
		AddUserId("Alice", "alice@example.com").
		New().
		GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	armored, err := key.Armor()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Classical private key (first 200 chars):\n%.200s...\n\n", armored)

	pubKey, err := key.ToPublic()
	if err != nil {
		log.Fatal(err)
	}
	pubArmored, err := pubKey.Armor()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Classical public key (first 200 chars):\n%.200s...\n\n", pubArmored)

	message := []byte("Hello from classical cryptography!")

	encHandle, err := pgp.Encryption().
		Recipient(key).
		New()
	if err != nil {
		log.Fatal(err)
	}

	ciphertext, err := encHandle.Encrypt(message)
	if err != nil {
		log.Fatal(err)
	}

	armoredMsg, err := ciphertext.ArmorBytes()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Encrypted message (first 200 chars):\n%.200s...\n\n", string(armoredMsg))

	decHandle, err := pgp.Decryption().
		DecryptionKey(key).
		New()
	if err != nil {
		log.Fatal(err)
	}

	decrypted, err := decHandle.Decrypt(armoredMsg, crypto.Armor)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Decrypted: %s\n\n", string(decrypted.Bytes()))
}

func pqcExample() {
	fmt.Println("--- 2. Post-Quantum Key Generation & Encryption (ML-KEM + ML-DSA) ---")

	pgp := crypto.PGPWithProfile(profile.PQC())

	key, err := pgp.KeyGeneration().
		AddUserId("Bob PQC", "bob-pqc@example.com").
		New().
		GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	armored, err := key.Armor()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("PQC private key (first 300 chars):\n%.300s...\n\n", armored)

	pubKey, err := key.ToPublic()
	if err != nil {
		log.Fatal(err)
	}
	pubArmored, err := pubKey.Armor()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("PQC public key (first 300 chars):\n%.300s...\n\n", pubArmored)

	fingerprint := key.GetFingerprint()
	fmt.Printf("Key fingerprint: %s\n\n", fingerprint)

	message := []byte("This message is protected against quantum computers!")

	encHandle, err := pgp.Encryption().
		Recipient(key).
		New()
	if err != nil {
		log.Fatal(err)
	}

	ciphertext, err := encHandle.Encrypt(message)
	if err != nil {
		log.Fatal(err)
	}

	armoredMsg, err := ciphertext.ArmorBytes()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("PQC encrypted message (first 300 chars):\n%.300s...\n\n", string(armoredMsg))

	decHandle, err := pgp.Decryption().
		DecryptionKey(key).
		New()
	if err != nil {
		log.Fatal(err)
	}

	decrypted, err := decHandle.Decrypt(armoredMsg, crypto.Armor)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("PQC Decrypted: %s\n\n", string(decrypted.Bytes()))
}

func signAndVerifyExample() {
	fmt.Println("--- 3. Inline Sign & Verify (PQC - ML-DSA + Ed25519) ---")

	pgp := crypto.PGPWithProfile(profile.PQC())

	key, err := pgp.KeyGeneration().
		AddUserId("Charlie", "charlie@example.com").
		New().
		GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	message := []byte("This document is legally binding.")

	signer, err := pgp.Sign().
		SigningKey(key).
		New()
	if err != nil {
		log.Fatal(err)
	}

	signed, err := signer.Sign(message, crypto.Armor)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Signed message (first 300 chars):\n%.300s...\n\n", string(signed))

	verifier, err := pgp.Verify().
		VerificationKey(key).
		New()
	if err != nil {
		log.Fatal(err)
	}

	verResult, err := verifier.VerifyInline(signed, crypto.Armor)
	if err != nil {
		log.Fatal(err)
	}

	if sigErr := verResult.SignatureError(); sigErr != nil {
		fmt.Printf("Signature INVALID: %v\n\n", sigErr)
	} else {
		fmt.Printf("Signature VALID! Message: %s\n\n", string(verResult.Bytes()))
	}
}

func detachedSignatureExample() {
	fmt.Println("--- 4. Detached Signature (PQC) ---")

	pgp := crypto.PGPWithProfile(profile.PQC())

	key, err := pgp.KeyGeneration().
		AddUserId("Dave", "dave@example.com").
		New().
		GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	document := []byte("Important file contents that need integrity verification.")

	signer, err := pgp.Sign().
		SigningKey(key).
		Detached().
		New()
	if err != nil {
		log.Fatal(err)
	}

	signature, err := signer.Sign(document, crypto.Armor)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Detached signature:\n%s\n", string(signature))

	verifier, err := pgp.Verify().
		VerificationKey(key).
		New()
	if err != nil {
		log.Fatal(err)
	}

	verResult, err := verifier.VerifyDetached(document, signature, crypto.Armor)
	if err != nil {
		log.Fatal(err)
	}

	if sigErr := verResult.SignatureError(); sigErr != nil {
		fmt.Printf("Detached signature INVALID: %v\n\n", sigErr)
	} else {
		fmt.Printf("Detached signature VALID!\n\n")
	}
}

func passwordEncryptionExample() {
	fmt.Println("--- 5. Password-Based Encryption (Symmetric, already quantum-safe) ---")

	pgp := crypto.PGP()

	password := []byte("my-strong-passphrase-2024")
	message := []byte("Secret message encrypted with a password (AES-256, quantum-resistant).")

	encHandle, err := pgp.Encryption().
		Password(password).
		New()
	if err != nil {
		log.Fatal(err)
	}

	ciphertext, err := encHandle.Encrypt(message)
	if err != nil {
		log.Fatal(err)
	}

	armoredMsg, err := ciphertext.ArmorBytes()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Password-encrypted message:\n%.200s...\n\n", string(armoredMsg))

	decHandle, err := pgp.Decryption().
		Password(password).
		New()
	if err != nil {
		log.Fatal(err)
	}

	decrypted, err := decHandle.Decrypt(armoredMsg, crypto.Armor)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Decrypted: %s\n\n", string(decrypted.Bytes()))
}
