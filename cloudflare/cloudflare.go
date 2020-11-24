package cloudflare

import (
	"fmt"
	"log"
	"strings"

	"github.com/cloudflare/cloudflare-go"
	"golang.org/x/net/publicsuffix"
)

type Cloudflare struct {
	logger *log.Logger
	api    *cloudflare.API
}

func NewCloudflareWithAPIToken(l *log.Logger, apiToken string) *Cloudflare {
	cf, err := cloudflare.NewWithAPIToken(apiToken)
	if err != nil {
		l.Fatalf("Error while initializing Cloudflare client\n%v", err)
	}

	return &Cloudflare{
		logger: l,
		api:    cf,
	}
}

// func NewCloudflareWithAPIKeyAndEmail(apiKey, email string) {
// 	// TODO
// }

func (cf *Cloudflare) ConfigureDomain(domain, content string) {
	cf.logger.Println("Configuring Cloudflare DNS...")
	baseDomain, err := publicsuffix.EffectiveTLDPlusOne(domain)
	if err != nil {
		cf.logger.Fatalf("Error while getting eTLD\n%v", err)
	}

	zoneID, err := cf.api.ZoneIDByName(baseDomain)
	if err != nil {
		cf.logger.Fatalf("Error while getting Cloudflare Zone ID for domain\n%v", err)
	}

	name := strings.Replace(domain, fmt.Sprintf(".%s", baseDomain), "", 1)

	if !cf.doesDNSRecordExist(domain) {
		_, err = cf.api.CreateDNSRecord(
			zoneID,
			cloudflare.DNSRecord{
				Type:    "CNAME",
				Name:    name,
				Content: content,
				Proxied: true,
			},
		)
		if err != nil {
			cf.logger.Fatalf("Error while creating Cloudflare DNS record\n%v", err)
		}
	}

	cf.logger.Println("Cloudflare DNS configuration completed")
}

func (cf *Cloudflare) doesDNSRecordExist(domain string) bool {
	baseDomain, err := publicsuffix.EffectiveTLDPlusOne(domain)
	if err != nil {
		cf.logger.Fatalf("Error while getting eTLD\n%v", err)
	}

	zoneID, err := cf.api.ZoneIDByName(baseDomain)
	if err != nil {
		cf.logger.Fatalf("Error while getting Cloudflare Zone ID for domain\n%v", err)
	}

	records, err := cf.api.DNSRecords(zoneID, cloudflare.DNSRecord{
		Type: "CNAME",
	})
	if err != nil {
		cf.logger.Fatalf("Error while checking if DNS CNAME record already exists\n%v", err)
	}

	for _, r := range records {
		if r.Name == domain {
			return true
		}
	}

	return false
}
