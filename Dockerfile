FROM debian:stable-slim
COPY dist/acmeproxy/acmeproxy /usr/bin/acmeproxy
RUN chmod +x /usr/bin/acmeproxy
RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*
ENTRYPOINT [ "/usr/bin/acmeproxy" ]
