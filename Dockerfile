FROM debian:stable-slim
COPY dist/acmeproxy/acmeproxy /usr/bin/acmeproxy
RUN chmod +x /usr/bin/acmeproxy
ENTRYPOINT [ "/usr/bin/acmeproxy" ]
