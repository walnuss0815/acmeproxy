package proxy

import (
	"strings"

	"github.com/walnuss0815/acmeproxy/v2/acme/provider"
)

type Proxy struct {
	AllowedDomains []string
	Provider       provider.Provider
}

func NewProxy(provider provider.Provider, allowedDomains []string) *Proxy {
	p := new(Proxy)
	p.AllowedDomains = allowedDomains
	p.Provider = provider

	return p
}

func (p *Proxy) CheckDomain(domain string) bool {
	for _, allowedDomain := range p.AllowedDomains {
		if strings.HasSuffix(domain, "."+allowedDomain+".") {
			return true
		}
	}

	return false
}

func (p *Proxy) Present(fqdn string, value string) error {
	err := p.Provider.CreateRecord(fqdn, value)
	return err
}

func (p *Proxy) Cleanup(fqdn string, value string) error {
	err := p.Provider.RemoveRecord(fqdn, value)
	return err
}
