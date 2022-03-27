package provider

type Provider interface {
	CreateRecord(fqdn, value string) error
	RemoveRecord(fqdn, value string) error
}
