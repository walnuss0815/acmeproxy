# ACMEproxy

Proxy server for ACME DNS01 challenges written in Go.

Works with the [httpreq](https://github.com/go-acme/lego/tree/master/providers/dns/httpreq) DNS challenge provider
in [lego](https://github.com/go-acme/lego) and with
the [acmeproxy](https://github.com/Neilpang/acme.sh/blob/dev/dnsapi/dns_acmeproxy.sh) provider in acme.sh.

## Why?

ACMEproxy was written to provide a way make it easier and safer to automatically issue
per-host [Let's Encrypt](https://letsencrypt.org) SSL certificates inside a larger network with many hosts. Especially
when these hosts aren't accessible from the outside, and they need to use the DNS challenges and require DNS API access.

The regular approach would be to run an ACME client on every host, which would also mean giving each hosts access to
the (full) DNS API. This is both hard to manage and a potential security risk.

As a solution ACMEproxy provides the following:

- Allow internal hosts to request ACME DNS challenges through a single host, without individual / full API access to the
  DNS provider
- Provide a single (acmeproxy) host that has access to the DNS credentials / API, limiting a possible attack surface
- Domain validation to only allow ACME DNS requests for specific domains

ACMEproxy was written to be run within an internal network, it's not recommended exposing your ACMEproxy host to the
outside world. Do so at your own risk.

Because ACMEproxy does not support https by itself, I recommend putting a reverse proxy like Traefik in front of
ACMEproxy.

# Build

## Prerequisite / WARNING

To use ACMEproxy as backend with providers from the `lego` package they need to implement a `CreateRecord`
/`RemoveRecord` method that takes an FQDN + acme value as input. The discussion if this should be practice is ongoing,
see [issue 720](https://github.com/go-acme/lego/issues/720). As an example take a look
at [PR #883](https://github.com/go-acme/lego/pull/883) of how this was implemented for the `transip` provider (don't
worry, it's not difficult).

Use the makefile to `make` the executables. Use `make install` to also install the executable to `/usr/local/bin`.

If you want to build a Debian package / installer, use `dch` to update the changelog and create your own package using `make debian`.


# Configure

You need to specify the relevant environment variables for the provider you've chosen. See
the [lego](https://github.com/go-acme/lego) documentation for options per provider.

| Name            | Argument  | Environment variable | Default | Example                                       |
| --------------- | --------- |----------------------|---------|-----------------------------------------------|
| Server port     | -port     | ACMEPROXY_PORT       | `9096`  | `9096`                                        |
| DNS provider    | -provider | ACMEPROXY_PROVIDER   |         | `cloudflare`                                  |
| Allowed domains | -domains  | ACMEPROXY_DOMAINS    |         | `internal1.example.com,internal2.example.com` |


# Usage

If you've configured ACMEproxy via the environment variables, you can just run `acmeproxy`.