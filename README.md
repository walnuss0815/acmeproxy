# acmeproxy
Proxy server for ACME DNS challenges written in Go.

Works with the [httpreq](https://github.com/go-acme/lego/tree/master/providers/dns/httpreq) DNS challenge provider in [lego](https://github.com/go-acme/lego) and with the [acmeproxy](https://github.com/Neilpang/acme.sh/blob/dev/dnsapi/dns_acmeproxy.sh) provider in acme.sh (currently in the dev branch).

## Why?
Acmeproxy was written to provide a way make it easier and safer to automatically issue per-host [Let's Encrypt](https://letsencrypt.org) SSL certificates inside a larger network with many different hosts. Especially when these hosts aren't accessible from the outside, and they need to use the DNS challenges and require DNS API access.

The regular approach would be to run an ACME client on every host, which would also mean giving each hosts access to the (full) DNS API. This is both hard to manage and a potential security risk.

As a solution Acmeproxy provides the following:
- Allow internal hosts to request ACME DNS challenges through a single host, without individual / full API access to the DNS provider
- Provide a single (acmeproxy) host that has access to the DNS credentials / API, limiting a possible attack surface
- Username/password or IP-based filtering for clients to prevent unauthorized access
- Domain validation to only allow ACME DNS requests for specific domains

If you're looking for other ways to validate internal certificates, take a look
at [autocertdelegate](https://github.com/bradfitz/autocertdelegate) which uses the tls-alpn-01 method.

Acmeproxy was written to be run within an internal network, it's not recommended exposing your Acmeproxy host to the
outside world. Do so at your own risk.

# Build

## Prerequisite / WARNING

to use acmeproxy as backend with providers from the `lego` package they need to implement a `CreateRecord`/`RemoveRecord` method that takes an FQDN + acme value as input. The discussion if this should be practice is on-going, see [issue 720](https://github.com/go-acme/lego/issues/720). As an example take a look at [PR #883](https://github.com/go-acme/lego/pull/883) of how this was implemented for the `transip` provider (don't worry, it's not difficult).

Use the makefile to `make` the executables. Use `make install` to also install the executable to `/usr/local/bin`.

If you want to build a Debian package / installer, use `dch` to update the changelog and create your own package using `make debian`.

# Configure

## Adjust configuration file
Copy `config.yml` to a directory (default: `/etc/acmeproxy`). See below for a configuration example using the `transip` provider. You need to specify the relevant environment variables for the provider you've chose. See the [lego](https://github.com/go-acme/lego) documentation for options per provider. Also see the examples below. If you want to provide proxies for multiple providers, start multiple instances on different hosts/ports (using different config files).

```
# Server configuration
server:
  port: 9096
  # htpasswd: "/etc/acmeproxy/htpasswd"
  accesslog: "/var/log/acmeproxy.log"

provider: "transip"

# Filter configuration
filter:
  ips:
    - "127.0.0.1"
    - "172.16.0.0/16"
  domains:
    - "example.com"

# Environment variables to be used with this provider
environment:
  TRANSIP_ACCOUNT_NAME: mdbraber
  TRANSIP_PRIVATE_KEY_PATH: /etc/acmeproxy/transip.key
  TRANSIP_POLLING_INTERVAL: 30
  TRANSIP_PROPAGATION_TIMEOUT: 600
```

## Authentication 
If you want to use client authentication (username/password), use following command: `htpasswd -c /etc/acmeproxy/htpasswd testuser` to create a new htpasswd file with user `testuser`.

If you want to use serverside IP based authentication set `allowed-ips` in the configfile (or set `--allowed-ips` on the commandline). You can use multiple IPs / nets in a CIDR notation, e.g. `127.0.0.1`, `172.16.0.0/16` or `192.168.10.0/24`.

# Usage

## Running acmeproxy in the foreground
If you've configured acmeproxy via the config file, you can just run `acmeproxy`. It will run in the foreground.

## Daemon mode
If you want to use acmeproxy as a daemon (in the background) use the `acmeproxy.service` in `debian/` as an example for systemd and copy it to `/etc/systemd/systemd` and enable it by `systemctl enable acmeproxy.service`. Be sure to check the `ExecStart` variable to see if it points to the right executable (`/usr/bin/acmeproxy` by default). Of course if you build `acmeproxy` as a Debian package the systemd service will be installed as part of the package.

## Showing systemd logs

If you run acmeproxy through systemd and use `log-forcecolors: true` and `log-forceformatting: true` - you can use `journalctl -xe -o cat -u acmeproxy.service` to see the original colored output with timestamps
