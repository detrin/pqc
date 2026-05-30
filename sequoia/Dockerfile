FROM rust:1.85-bookworm AS builder

RUN apt-get update && apt-get install -y \
    pkg-config \
    capnproto \
    clang \
    libclang-dev \
    perl \
    make \
    && rm -rf /var/lib/apt/lists/*

# OpenSSL 3.5 required for ML-KEM/ML-DSA support
RUN curl -L -o openssl.tar.gz https://github.com/openssl/openssl/releases/download/openssl-3.5.0/openssl-3.5.0.tar.gz \
    && echo "344d0a79f1a9b08029b0744e2cc401a43f9c90acd1044d09a530b4885a8e9fc0  openssl.tar.gz" | sha256sum -c \
    && tar xzf openssl.tar.gz && rm openssl.tar.gz \
    && cd openssl-3.5.0 \
    && ./Configure --prefix=/opt/openssl --openssldir=/opt/openssl/ssl --libdir=lib \
    && make -j$(nproc) \
    && make install_sw \
    && cd .. && rm -rf openssl-3.5.0

ENV OPENSSL_DIR=/opt/openssl
ENV PKG_CONFIG_PATH=/opt/openssl/lib/pkgconfig
ENV BINDGEN_EXTRA_CLANG_ARGS="-I/opt/openssl/include"

RUN cargo install sequoia-sq --version 1.4.0-pqc.1 \
    --locked --no-default-features --features crypto-openssl

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y \
    libsqlite3-0 \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /opt/openssl/lib/libssl.so.3 /usr/local/lib/
COPY --from=builder /opt/openssl/lib/libcrypto.so.3 /usr/local/lib/
RUN ldconfig
COPY --from=builder /usr/local/cargo/bin/sq /usr/local/bin/sq

WORKDIR /data
ENTRYPOINT ["sq"]
CMD ["version"]
