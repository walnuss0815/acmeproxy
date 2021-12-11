FROM debian:stable-slim
RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*
COPY dist/acmeproxy/acmeproxy /usr/bin/acmeproxy
RUN chmod +x /usr/bin/acmeproxy
ENTRYPOINT [ "/usr/bin/acmeproxy" ]
