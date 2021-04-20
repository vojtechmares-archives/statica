package cloudflare

import (
	"context"
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
			context.Background(),
			zoneID,
			cloudflare.DNSRecord{
				Type:    "CNAME",
				Name:    name,
				Content: content,
				Proxied: pTrue(),
			},
		)
		if err != nil {
			cf.logger.Fatalf("Error while creating Cloudflare DNS record\n%v", err)
		}
	}

	cf.logger.Println("Cloudflare DNS configuration completed")
}

func (cf *Cloudflare) DestroyDomain(domain string) {
	cf.logger.Println("Deleting Cloudflare DNS record...")
	baseDomain, err := publicsuffix.EffectiveTLDPlusOne(domain)
	if err != nil {
		cf.logger.Fatalf("Error while getting eTLD\n%v", err)
	}

	zoneID, err := cf.api.ZoneIDByName(baseDomain)
	if err != nil {
		cf.logger.Fatalf("Error while getting Cloudflare Zone ID for domain\n%v", err)
	}

	records, err := cf.api.DNSRecords(
		context.Background(),
		zoneID,
		cloudflare.DNSRecord{
			Type: "CNAME",
		})
	if err != nil {
		cf.logger.Fatalf("Error while getting DNS records from Cloudflare\n%v", err)
	}

	var recordID string

	for _, r := range records {
		if r.Name == domain {
			recordID = r.ID
		}
	}

	err = cf.api.DeleteDNSRecord(context.Background(), zoneID, recordID)
	if err != nil {
		cf.logger.Fatalf("Error while deleting Cloudflare DNS record\n%v", err)
	}

	cf.logger.Println("Cloudflare DNS record deleted")
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

	records, err := cf.api.DNSRecords(
		context.Background(),
		zoneID,
		cloudflare.DNSRecord{
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

func pTrue() *bool {
	t := true
	return &t
}
