package main

import (
	"context"
	"fmt"
	"os"

	"github.com/digitalocean/godo"
	"github.com/thenets/do-kyoka/helper"
)

func main() {
	// Read env vars
	firewallName := os.Getenv("FIREWALL_NAME")
	if firewallName == "" {
		panic("[ERROR] 'FIREWALL_NAME' not set in config file!")
	}
	tagName := os.Getenv("FIREWALL_TAG")
	if tagName == "" {
		panic("[ERROR] 'FIREWALL_TAG' not set in config file!")
	}
	apiToken := os.Getenv("DO_API_TOKEN")
	if tagName == "" {
		panic("[ERROR] 'DO_API_TOKEN' not set in config file!")
	}
	
	// DEBUG
	fmt.Println("firewallName", firewallName)
	fmt.Println("tagName", tagName)
	fmt.Println("apiToken", apiToken)
	return

	// Firewall: allow my current IP
	client := godo.NewFromToken(apiToken)
	ctx := context.TODO()
	firewall, err := helper.FirewallAllowMyCurrentIp(ctx, client, firewallName, tagName)
	if err != nil {
		panic(nil)
	}

	// Output Firewall information
	fmt.Println(firewall)
}
