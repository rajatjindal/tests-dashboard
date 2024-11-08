#FROM ghcr.io/fermyon/spin:v3.0.0-rc.1-distroless
FROM rajatjindal/spin:2906@sha256:91bcca98b9bdbb42757bb6266b143bf2639d9af42fed3e1e7d63566c06dfcbbd

WORKDIR /app

COPY spin.toml spin.toml
COPY backend/main.wasm backend/main.wasm
COPY redirector/redirect.wasm redirector/redirect.wasm
COPY runtimeconfig.toml runtimeconfig.toml
COPY ui/.output/public ui/.output/public

ENTRYPOINT ["spin", "up"]
