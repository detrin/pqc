# GopenPGP PQC (pqcrypt)

A minimal PQC-capable PGP CLI tool built on [Proton's GopenPGP](https://github.com/ProtonMail/gopenpgp) library.

## Usage

```bash
docker build -t pqcrypt .
```

```bash
# Generate PQC key pair
docker run --rm -v $(pwd):/data pqcrypt keygen \
  --name "Alice" --email alice@example.com --pqc --output /data/alice

# Encrypt
docker run --rm -v $(pwd):/data pqcrypt encrypt \
  --recipient /data/alice.pub.asc --input /data/secret.txt

# Decrypt
docker run --rm -v $(pwd):/data pqcrypt decrypt \
  --key /data/alice.key.asc --input /data/secret.txt.pgp

# Sign
docker run --rm -v $(pwd):/data pqcrypt sign \
  --key /data/alice.key.asc --input /data/document.txt

# Verify
docker run --rm -v $(pwd):/data pqcrypt verify \
  --pubkey /data/alice.pub.asc --input /data/document.txt

# Inspect key
docker run --rm -v $(pwd):/data pqcrypt inspect --key /data/alice.pub.asc
```

## Build Locally

```bash
go build -o pqcrypt ./cmd/pqcrypt/
```

## Library

`main.go` contains standalone examples of GopenPGP usage:
- Classical key generation (Ed25519/X25519)
- PQC key generation (ML-DSA-65+Ed25519 / ML-KEM-768+X25519)
- Encrypt/decrypt with PQC
- Sign/verify with PQC
- Password-based encryption

## Dependencies

- `github.com/ProtonMail/gopenpgp/v3@v3.4.1-proton`
- `github.com/ProtonMail/go-crypto@v1.4.1-proton`

The `-proton` tagged releases include PQC support from `draft-ietf-openpgp-pqc`.
