# Sequoia PQC

Post-quantum `sq` CLI in a Docker container. No build dependencies needed.

## Usage

```bash
docker build -t sequoia-pqc .
```

```bash
# Show version
docker run --rm sequoia-pqc

# Generate PQC key (ML-DSA-65+Ed25519 / ML-KEM-768+X25519)
docker run --rm -v $(pwd):/data sequoia-pqc key generate \
  --own-key --name "Alice" --email alice@example.com \
  --without-password --profile rfc9580 \
  --cipher-suite mldsa65-ed25519 \
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

## Build from Source (macOS)

```bash
brew install openssl@3 capnp

BINDGEN_EXTRA_CLANG_ARGS="-I$(brew --prefix openssl@3)/include" \
OPENSSL_DIR=$(brew --prefix openssl@3) \
C_INCLUDE_PATH=$(brew --prefix openssl@3)/include \
LIBRARY_PATH=$(brew --prefix openssl@3)/lib \
  cargo install sequoia-sq --version 1.4.0-pqc.1 \
    --locked --no-default-features --features crypto-openssl
```

## What's Inside

- **sq 1.4.0-pqc.1** (Sequoia PGP pre-release with PQC)
- **OpenSSL 3.5.0** (required for ML-KEM/ML-DSA algorithm support)
- Based on `debian:bookworm-slim`
