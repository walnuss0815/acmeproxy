package provider

import (
	"context"
	"fmt"
	"sync"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/go-acme/lego/challenge/dns01"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ProviderCloudflare struct {
	api *cloudflare.API

	recordIDs   map[string]string
	recordIDsMu sync.Mutex
}

func NewDefaultProviderCloudflare() *ProviderCloudflare {
	provider := new(ProviderCloudflare)

	useApiToken := viper.IsSet("provider.api_token")
	if useApiToken {
		apiToken := viper.GetString("provider.api_token")

		api, err := cloudflare.NewWithAPIToken(apiToken)
		if err != nil {
			log.Fatal(err)
		}

		provider.api = api
	} else {
		apiKey := viper.GetString("provider.api_key")
		apiEmail := viper.GetString("provider.api_email")

		api, err := cloudflare.New(apiKey, apiEmail)
		if err != nil {
			log.Fatal(err)
		}

		provider.api = api
	}

	provider.recordIDs = make(map[string]string)

	return provider
}

func (p *ProviderCloudflare) CreateRecord(fqdn string, value string) error {
	authZone, err := dns01.FindZoneByFqdn(fqdn)
	if err != nil {
		return fmt.Errorf("cloudflare: %w", err)
	}

	zoneID, err := p.api.ZoneIDByName(authZone)
	if err != nil {
		return fmt.Errorf("cloudflare: failed to find zone %s: %w", authZone, err)
	}

	dnsRecord := cloudflare.DNSRecord{
		Type:    "TXT",
		Name:    dns01.UnFqdn(fqdn),
		Content: value,
		TTL:     120,
	}

	response, err := p.api.CreateDNSRecord(context.Background(), zoneID, dnsRecord)
	if err != nil {
		return fmt.Errorf("cloudflare: failed to create TXT record: %w", err)
	}

	if !response.Success {
		return fmt.Errorf("cloudflare: failed to create TXT record: %+v %+v", response.Errors, response.Messages)
	}

	p.recordIDsMu.Lock()
	p.recordIDs[fqdn] = response.Result.ID
	p.recordIDsMu.Unlock()

	return nil
}

func (p *ProviderCloudflare) RemoveRecord(fqdn string, value string) error {
	authZone, err := dns01.FindZoneByFqdn(fqdn)
	if err != nil {
		return fmt.Errorf("cloudflare: %w", err)
	}

	zoneID, err := p.api.ZoneIDByName(authZone)
	if err != nil {
		return fmt.Errorf("cloudflare: failed to find zone %s: %w", authZone, err)
	}

	p.recordIDsMu.Lock()
	recordID, ok := p.recordIDs[fqdn]
	p.recordIDsMu.Unlock()
	if !ok {
		return fmt.Errorf("cloudflare: unknown record ID for '%s'", fqdn)
	}

	err = p.api.DeleteDNSRecord(context.Background(), zoneID, recordID)
	if err != nil {
		log.Printf("cloudflare: failed to delete TXT record: %w", err)
	}

	p.recordIDsMu.Lock()
	delete(p.recordIDs, fqdn)
	p.recordIDsMu.Unlock()

	return nil
}
