package helper

import (
	"context"

	"github.com/digitalocean/godo"
)

func ProjectNameList(ctx context.Context, client *godo.Client) ([]string, error) {
	var project_names []string

	opt := &godo.ListOptions{
		PerPage: 10,
		Page:    1,
	}
	projects, _, err := client.Projects.List(ctx, opt)

	for _, project := range projects {
		project_names = append(project_names, project.Name)
	}

	return project_names, err
}

func GetFirewallList(ctx context.Context, client *godo.Client) ([]godo.Firewall, error) {
	var firewallList []godo.Firewall
	var err error

	opt := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}
	firewallList, _, err = client.Firewalls.List(ctx, opt)

	return firewallList, err
}

func HasFirewall(ctx context.Context, client *godo.Client, firewallName string) bool {
	firewalls, err := GetFirewallList(ctx, client)
	if err != nil {
		return false
	}

	// Search
	for _, firewall := range firewalls {
		if firewall.Name == firewallName {
			return true
		}
	}

	return false
}

func GetFirewallByName(ctx context.Context, client *godo.Client, firewallName string) (*godo.Firewall, error) {
	var err error
	var firewallNull *godo.Firewall

	firewalls, err := GetFirewallList(ctx, client)
	if err != nil {
		return firewallNull, err
	}

	for _, firewall := range firewalls {
		if firewall.Name == firewallName {
			return &firewall, err
		}
	}

	return firewallNull, err
}

func FirewallAllowMyCurrentIp(ctx context.Context, client *godo.Client, firewallName string, tagName string) (*godo.Firewall, error) {
	var firewall *godo.Firewall
	var err error

	// Get current public IP
	currentIp, err := GetMyPublicIp()
	if err != nil {
		panic("[ERROR] Can't get public IP address!")
	}
	// fmt.Println(currentIp)

	// Create tag
	_, _, err = client.Tags.Get(ctx, "awesome")
	if err != nil {
		createRequest := &godo.TagCreateRequest{
			Name: tagName,
		}
		_, _, err = client.Tags.Create(ctx, createRequest)
		if err != nil {
			return firewall, err
		}
	}

	// Check if Firewall exist
	if !HasFirewall(ctx, client, firewallName) {
		// Create Firewall if not exist
		createRequest := &godo.FirewallRequest{
			Name: firewallName,
			InboundRules: []godo.InboundRule{
				{
					Protocol:  "tcp",
					PortRange: "22",
					Sources: &godo.Sources{
						Addresses: []string{"127.0.0.1/24"},
					},
				},
			},
		}
		_, _, err := client.Firewalls.Create(ctx, createRequest)
		if err != nil {
			return firewall, err
		}
	}

	// Get Firewall
	firewall, err = GetFirewallByName(ctx, client, firewallName)
	if err != nil {
		return firewall, err
	}

	// Update firewall with new rules
	updateRequest := &godo.FirewallRequest{
		Name: firewallName,
		InboundRules: []godo.InboundRule{
			{
				Protocol:  "tcp",
				PortRange: "all",
				Sources: &godo.Sources{
					Addresses: []string{currentIp},
				},
			},
		},
		Tags: []string{tagName},
	}
	firewall, _, err = client.Firewalls.Update(ctx, firewall.ID, updateRequest)
	if err != nil {
		return firewall, err
	}

	return firewall, err
}
