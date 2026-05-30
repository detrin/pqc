#!/bin/bash
set -euo pipefail

# ============================================================================
# Sequoia PGP - Post-Quantum Cryptography Demo
# ============================================================================
#
# This script demonstrates post-quantum cryptography using Sequoia PGP's
# sq CLI tool with ML-DSA-65+Ed25519 (signing) and ML-KEM-768+X25519 (encryption).
#
# Prerequisites:
#   1. Rust toolchain: curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
#   2. OpenSSL 3.x:   brew install openssl@3
#   3. Cap'n Proto:    brew install capnp
#
# Build Sequoia PQC (one-time, takes ~1-2 minutes):
#   BINDGEN_EXTRA_CLANG_ARGS="-I$(brew --prefix openssl@3)/include" \
#   OPENSSL_DIR=$(brew --prefix openssl@3) \
#   C_INCLUDE_PATH=$(brew --prefix openssl@3)/include \
#   LIBRARY_PATH=$(brew --prefix openssl@3)/lib \
#     cargo install sequoia-sq --version 1.4.0-pqc.1 \
#       --locked --no-default-features --features crypto-openssl
#
# After build, sq is at: ~/.cargo/bin/sq
# ============================================================================

SQ="${SQ:-$(command -v sq 2>/dev/null || echo "$HOME/.cargo/bin/sq")}"
WORKDIR="$(cd "$(dirname "$0")" && pwd)"
export SEQUOIA_HOME="$WORKDIR/.sequoia-home"

mkdir -p "$WORKDIR/keys" "$WORKDIR/messages" "$WORKDIR/signatures" "$SEQUOIA_HOME"

echo "Using: $($SQ version | head -1)"
echo "Working directory: $WORKDIR"
echo

# --------------------------------------------------------------------------
# Step 1: Generate a post-quantum key pair
# --------------------------------------------------------------------------
echo "=== Step 1: Generate PQC key pair (ML-DSA-65+Ed25519 / ML-KEM-768+X25519) ==="
$SQ --home="$SEQUOIA_HOME" key generate \
  --own-key \
  --name "Alice" \
  --email alice@example.com \
  --without-password \
  --profile rfc9580 \
  --cipher-suite mldsa65-ed25519 \
  --output "$WORKDIR/keys/alice.secret.pgp" \
  --rev-cert "$WORKDIR/keys/alice.rev"
echo

# --------------------------------------------------------------------------
# Step 2: Extract public certificate (strip secret key material)
# --------------------------------------------------------------------------
echo "=== Step 2: Extract public certificate ==="
$SQ --home="$SEQUOIA_HOME" key delete \
  --cert-file="$WORKDIR/keys/alice.secret.pgp" \
  --output="$WORKDIR/keys/alice.pub.pgp"
echo "Public key: $WORKDIR/keys/alice.pub.pgp"
echo

# --------------------------------------------------------------------------
# Step 3: Encrypt a message using the PQC public key
# --------------------------------------------------------------------------
echo "=== Step 3: Encrypt a message ==="
echo "This message is protected by ML-KEM-768+X25519 post-quantum encryption." \
  > "$WORKDIR/messages/plaintext.txt"

$SQ --home="$SEQUOIA_HOME" encrypt \
  --for-file "$WORKDIR/keys/alice.pub.pgp" \
  --without-signature \
  --output "$WORKDIR/messages/encrypted.pgp" \
  "$WORKDIR/messages/plaintext.txt"
echo

# --------------------------------------------------------------------------
# Step 4: Decrypt the message using the private key
# --------------------------------------------------------------------------
echo "=== Step 4: Decrypt the message ==="
$SQ --home="$SEQUOIA_HOME" decrypt \
  --recipient-file "$WORKDIR/keys/alice.secret.pgp" \
  --output "$WORKDIR/messages/decrypted.txt" \
  "$WORKDIR/messages/encrypted.pgp"
echo "Decrypted content:"
cat "$WORKDIR/messages/decrypted.txt"
echo

# --------------------------------------------------------------------------
# Step 5: Sign a document with ML-DSA-65+Ed25519 (detached signature)
# --------------------------------------------------------------------------
echo "=== Step 5: Sign a document ==="
echo "I hereby declare this document authentic." > "$WORKDIR/messages/document.txt"

$SQ --home="$SEQUOIA_HOME" sign \
  --signer-file "$WORKDIR/keys/alice.secret.pgp" \
  --signature-file "$WORKDIR/signatures/document.txt.sig" \
  "$WORKDIR/messages/document.txt"
echo "Signature: $WORKDIR/signatures/document.txt.sig"
echo

# --------------------------------------------------------------------------
# Step 6: Verify the signature
# --------------------------------------------------------------------------
echo "=== Step 6: Verify the signature ==="
$SQ --home="$SEQUOIA_HOME" verify \
  --signer-file "$WORKDIR/keys/alice.pub.pgp" \
  --signature-file "$WORKDIR/signatures/document.txt.sig" \
  "$WORKDIR/messages/document.txt"
echo

# --------------------------------------------------------------------------
# Step 7: Inspect key information
# --------------------------------------------------------------------------
echo "=== Step 7: Inspect key ==="
$SQ --home="$SEQUOIA_HOME" inspect "$WORKDIR/keys/alice.pub.pgp"
echo

# --------------------------------------------------------------------------
# Summary
# --------------------------------------------------------------------------
echo "=== Summary ==="
echo "All operations completed successfully using post-quantum cryptography."
echo
echo "Algorithms used:"
echo "  Signing:    ML-DSA-65 + Ed25519 (hybrid)"
echo "  Encryption: ML-KEM-768 + X25519 (hybrid)"
echo "  Symmetric:  AES-256/OCB"
echo
echo "Files:"
echo "  Private key:     keys/alice.secret.pgp"
echo "  Public key:      keys/alice.pub.pgp"
echo "  Revocation cert: keys/alice.rev"
echo "  Encrypted msg:   messages/encrypted.pgp"
echo "  Decrypted msg:   messages/decrypted.txt"
echo "  Document:        messages/document.txt"
echo "  Signature:       signatures/document.txt.sig"
echo
echo "File sizes (PQC overhead):"
wc -c "$WORKDIR/keys/alice.secret.pgp" "$WORKDIR/keys/alice.pub.pgp" \
      "$WORKDIR/messages/encrypted.pgp" "$WORKDIR/signatures/document.txt.sig" \
  2>/dev/null | grep -v total
echo
echo "Standard: draft-ietf-openpgp-pqc (implemented via sequoia-openpgp 2.2.0-pqc.1)"
