FROM alpine:latest@sha256:1e42bbe2508154c9126d48c2b8a75420c3544343bf86fd041fb7527e017a4b4a

ARG TARGETARCH

LABEL \
    org.opencontainers.image.title="eth-peer-manager" \
    org.opencontainers.image.description="Eth peer manager service" \
    org.opencontainers.image.source="https://github.com/jer117/eth-peer-manager" \
    org.opencontainers.image.vendor="Blockdaemon Inc." \
    org.opencontainers.image.version="${VERSION:-untagged}"

ADD ./build/eth-peer-manager-linux-${TARGETARCH} /eth-peer-manager

ENTRYPOINT ["/eth-peer-manager"]
