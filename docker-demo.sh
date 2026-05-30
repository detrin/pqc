#!/bin/bash
set -euo pipefail

echo "============================================================"
echo " Post-Quantum Cryptography in PGP - Linux Demo"
echo "============================================================"
echo
echo "Tools:"
sq version
echo
pqcrypt --help | head -1
echo

export SEQUOIA_HOME=/demo/.sequoia-home
mkdir -p /demo/keys /demo/messages /demo/signatures "$SEQUOIA_HOME"

# -------------------------------------------------------------------------
echo "=== 1. Generate PQC key with Sequoia (ML-DSA-65+Ed25519) ==="
# -------------------------------------------------------------------------
sq --home="$SEQUOIA_HOME" key generate \
  --own-key \
  --name "Linux PQC User" \
  --email pqc@linux-demo.org \
  --without-password \
  --profile rfc9580 \
  --cipher-suite mldsa65-ed25519 \
  --output /demo/keys/user.secret.pgp \
  --rev-cert /demo/keys/user.rev

sq --home="$SEQUOIA_HOME" key delete \
  --cert-file=/demo/keys/user.secret.pgp \
  --output=/demo/keys/user.pub.pgp
echo

# -------------------------------------------------------------------------
echo "=== 2. Encrypt with Sequoia, Decrypt with pqcrypt ==="
# -------------------------------------------------------------------------
echo "Secret message from Linux container!" > /demo/messages/secret.txt

sq --home="$SEQUOIA_HOME" encrypt \
  --for-file /demo/keys/user.pub.pgp \
  --without-signature \
  --output /demo/messages/secret.pgp \
  /demo/messages/secret.txt

pqcrypt decrypt \
  --key /demo/keys/user.secret.pgp \
  --input /demo/messages/secret.pgp \
  --output /demo/messages/decrypted.txt

echo "Original:  $(cat /demo/messages/secret.txt)"
echo "Decrypted: $(cat /demo/messages/decrypted.txt)"
echo

# -------------------------------------------------------------------------
echo "=== 3. Encrypt with pqcrypt, Decrypt with Sequoia ==="
# -------------------------------------------------------------------------
echo "Another secret, encrypted by pqcrypt." > /demo/messages/from-go.txt

pqcrypt encrypt \
  --recipient /demo/keys/user.pub.pgp \
  --input /demo/messages/from-go.txt \
  --output /demo/messages/from-go.pgp

sq --home="$SEQUOIA_HOME" decrypt \
  --recipient-file /demo/keys/user.secret.pgp \
  --output /demo/messages/from-go-dec.txt \
  /demo/messages/from-go.pgp

echo "Decrypted: $(cat /demo/messages/from-go-dec.txt)"
echo

# -------------------------------------------------------------------------
echo "=== 4. Sign with pqcrypt, Verify with Sequoia ==="
# -------------------------------------------------------------------------
echo "Document signed on Linux." > /demo/messages/doc.txt

pqcrypt sign --key /demo/keys/user.secret.pgp --input /demo/messages/doc.txt

sq --home="$SEQUOIA_HOME" verify \
  --signer-file /demo/keys/user.pub.pgp \
  --signature-file /demo/messages/doc.txt.sig \
  /demo/messages/doc.txt
echo

# -------------------------------------------------------------------------
echo "=== 5. Sign with Sequoia, Verify with pqcrypt ==="
# -------------------------------------------------------------------------
sq --home="$SEQUOIA_HOME" sign \
  --signer-file /demo/keys/user.secret.pgp \
  --signature-file /demo/signatures/doc-sq.sig \
  /demo/messages/doc.txt

pqcrypt verify \
  --pubkey /demo/keys/user.pub.pgp \
  --input /demo/messages/doc.txt \
  --signature /demo/signatures/doc-sq.sig
echo

# -------------------------------------------------------------------------
echo "=== 6. Generate classical key with pqcrypt for comparison ==="
# -------------------------------------------------------------------------
pqcrypt keygen --name "Classical User" --email classical@demo.org --output /demo/keys/classical
echo

# -------------------------------------------------------------------------
echo "=== 7. File size comparison ==="
# -------------------------------------------------------------------------
echo
echo "PQC key (ML-DSA-65+Ed25519 / ML-KEM-768+X25519):"
wc -c /demo/keys/user.pub.pgp
echo "Classical key (Ed25519 / X25519):"
wc -c /demo/keys/classical.pub.asc
echo
echo "PQC encrypted message:"
wc -c /demo/messages/secret.pgp
echo "PQC signature:"
wc -c /demo/messages/doc.txt.sig
echo

# -------------------------------------------------------------------------
echo "============================================================"
echo " All tests passed! PQC works on Linux ($(uname -m))."
echo " Both tools interoperate: Sequoia (Rust) <-> pqcrypt (Go)"
echo "============================================================"
