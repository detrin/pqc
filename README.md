# PQC — Post-Quantum Cryptography in PGP

Practical post-quantum encryption and signing using OpenPGP. Two implementations, both interoperable, each with a ready-to-use Docker image.

## Why?

NIST standardized ML-KEM and ML-DSA in 2024. The IETF OpenPGP working group is ratifying `draft-ietf-openpgp-pqc` to bring these algorithms into PGP. But there's no easy way to try it — the tools require OpenSSL 3.5+ (which no distro ships), building Rust from source, or tracking down pre-release Go modules.

This repo packages everything into Docker images so you can generate post-quantum keys, encrypt files, and sign documents in one command.

> **Warning:** This uses pre-release implementations of a draft standard. The algorithms (ML-KEM-768, ML-DSA-65) are NIST-standardized, but the OpenPGP integration is not yet an RFC. Do not use for production secrets without understanding the risks. Keys may not be compatible with future stable releases.

## Quick Start

### Option A: Pull from GitHub Container Registry (no clone needed)

```bash
docker pull ghcr.io/detrin/pqc/sequoia:latest
docker pull ghcr.io/detrin/pqc/gopenpgp:latest
```

Then use directly:

```bash
# Sequoia PQC
docker run --rm -v $(pwd):/data ghcr.io/detrin/pqc/sequoia key generate \
  --own-key --name "Alice" --email alice@example.com \
  --without-password --profile rfc9580 --cipher-suite mldsa65-ed25519 \
  --output /data/alice.pgp --rev-cert /data/alice.rev

# GopenPGP (pqcrypt)
docker run --rm -v $(pwd):/data ghcr.io/detrin/pqc/gopenpgp keygen \
  --name "Alice" --email alice@example.com --pqc --output /data/alice
```

### Option B: Clone and build locally

```bash
git clone https://github.com/detrin/pqc.git
cd pqc

# Build Sequoia PQC image
docker build -t sequoia-pqc ./sequoia/

# Build GopenPGP image
docker build -t pqcrypt ./gopenpgp/
```

## Usage: Sequoia PQC

The Sequoia image wraps the `sq` CLI with post-quantum support:

```bash
# Generate PQC key pair (ML-DSA-65+Ed25519 / ML-KEM-768+X25519)
docker run --rm -v $(pwd):/data sequoia-pqc key generate \
  --own-key --name "Alice" --email alice@example.com \
  --without-password --profile rfc9580 --cipher-suite mldsa65-ed25519 \
  --output /data/alice.pgp --rev-cert /data/alice.rev

# Encrypt
docker run --rm -v $(pwd):/data sequoia-pqc encrypt \
  --for-file /data/alice.pgp --without-signature \
  --output /data/secret.pgp /data/secret.txt

# Decrypt
docker run --rm -v $(pwd):/data sequoia-pqc decrypt \
  --recipient-file /data/alice.pgp \
  --output /data/decrypted.txt /data/secret.pgp

# Sign
docker run --rm -v $(pwd):/data sequoia-pqc sign \
  --signer-file /data/alice.pgp \
  --signature-file /data/doc.sig /data/doc.txt

# Verify
docker run --rm -v $(pwd):/data sequoia-pqc verify \
  --signer-file /data/alice.pgp \
  --signature-file /data/doc.sig /data/doc.txt
```

## Usage: GopenPGP (pqcrypt)

The GopenPGP image wraps `pqcrypt`, a simpler CLI built on Proton Mail's library:

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
  --key /data/alice.key.asc --input /data/doc.txt

# Verify
docker run --rm -v $(pwd):/data pqcrypt verify \
  --pubkey /data/alice.pub.asc --input /data/doc.txt
```

## Algorithms

| Purpose | Algorithm | Security Level |
|---------|-----------|----------------|
| Signing | ML-DSA-65 + Ed25519 (hybrid) | NIST Level 3 (~AES-192) |
| Encryption | ML-KEM-768 + X25519 (hybrid) | NIST Level 3 (~AES-192) |
| Symmetric | AES-256/OCB | Already quantum-safe |

Hybrid means both algorithms must be broken to compromise the key — quantum resistance + classical fallback.

## Interoperability

Both tools implement `draft-ietf-openpgp-pqc` and are fully interoperable:

| Test | Result |
|------|--------|
| Sequoia encrypts → GopenPGP decrypts | ✓ |
| GopenPGP encrypts → Sequoia decrypts | ✓ |
| Sequoia signs → GopenPGP verifies | ✓ |
| GopenPGP signs → Sequoia verifies | ✓ |

Run the full interop test:
```bash
docker build -t pqc-interop .
docker run --rm pqc-interop
```

## Project Structure

```
.
├── sequoia/                # Sequoia PGP (Rust)
│   ├── Dockerfile          # sq 1.4.0-pqc.1 + OpenSSL 3.5
│   ├── reproduce.sh        # Full demo script
│   └── README.md
├── gopenpgp/               # Proton GopenPGP (Go)
│   ├── Dockerfile          # pqcrypt CLI
│   ├── cmd/pqcrypt/        # CLI source
│   ├── main.go             # Library usage examples
│   └── README.md
├── .github/workflows/      # Auto-builds and pushes images to GHCR
├── Dockerfile              # Combined interop test
└── docker-demo.sh
```

## Build from Source (no Docker)

### Sequoia PQC (macOS)

```bash
brew install openssl@3 capnp

BINDGEN_EXTRA_CLANG_ARGS="-I$(brew --prefix openssl@3)/include" \
OPENSSL_DIR=$(brew --prefix openssl@3) \
C_INCLUDE_PATH=$(brew --prefix openssl@3)/include \
LIBRARY_PATH=$(brew --prefix openssl@3)/lib \
  cargo install sequoia-sq --version 1.4.0-pqc.1 \
    --locked --no-default-features --features crypto-openssl
```

### pqcrypt (any platform with Go)

```bash
cd gopenpgp/
go build -o pqcrypt ./cmd/pqcrypt/
```

## Standard

Based on [draft-ietf-openpgp-pqc](https://datatracker.ietf.org/doc/draft-ietf-openpgp-pqc/) (nearing IETF ratification).
