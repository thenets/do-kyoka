package main

import (
	"context"
	"fmt"
	"os"

	"github.com/digitalocean/godo"
	"github.com/jasonlvhit/gocron"
	log "github.com/sirupsen/logrus"
	"github.com/thenets/do-kyoka/helper"
)

func updateFirewall() {
	// Read env vars
	firewallName := os.Getenv("FIREWALL_NAME")
	if firewallName == "" {
		log.Panic("'FIREWALL_NAME' not set in config file!")
	}
	tagName := os.Getenv("FIREWALL_TAG")
	if tagName == "" {
		log.Panic("'FIREWALL_TAG' not set in config file!")
	}
	apiToken := os.Getenv("DO_API_TOKEN")
	if tagName == "" {
		log.Panic("'DO_API_TOKEN' not set in config file!")
	}

	// Firewall: allow my current IP
	client := godo.NewFromToken(apiToken)
	ctx := context.TODO()
	firewall, err := helper.FirewallAllowMyCurrentIp(ctx, client, firewallName, tagName)
	if err != nil {
		// Output Firewall information
		log.WithFields(log.Fields{
			"ID":            firewall.ID,
			// "Name":          firewall.Name,
			// "Status":        firewall.Status,
			"InboundRules":  firewall.InboundRules,
			"OutboundRules": firewall.OutboundRules,
		}).Error(fmt.Sprintf("firewall '%s' update %s", firewall.Name, firewall.Status))
	}

	// Output Firewall information
	log.WithFields(log.Fields{
		"ID":            firewall.ID,
		// "Name":          firewall.Name,
		// "Status":        firewall.Status,
		// "InboundRules":  firewall.InboundRules,
		// "OutboundRules": firewall.OutboundRules,
	}).Info(fmt.Sprintf("firewall '%s' update %s", firewall.Name, firewall.Status))
}

func main() {
	// Run updateFirewall once
	updateFirewall()

	// Add updateFirewall to scheduler
	gocron.Every(30).Minutes().Do(updateFirewall)

	// Start scheduler
	<-gocron.Start()
}
