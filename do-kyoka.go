package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/jasonlvhit/gocron"
	"github.com/sirupsen/logrus"

	log "github.com/thenets/go-sentry-logger"

	"github.com/digitalocean/godo"
)

func GetMyPublicIp() (string, error) {
	// Use CloudFlare instead
	// https://cloudflare.com/cdn-cgi/trace
	// https://gist.github.com/ankanch/8c8ec5aaf374039504946e7e2b2cdf7f
	url := "https://api.ipify.org?format=text"
	// fmt.Println("Getting IP address from 'ipify'...")

	resp, err := http.Get(url)
	if err != nil {
		return "nil", err
	}
	defer resp.Body.Close()

	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "nil", err
	}

	log.Info("Public IP identified: " + string(ip))
	log.CaptureMessage("yay! " + string(ip))

	return string(ip), nil
}

// IsPublicIP check if an IP is in public range
// Example: IsPublicIP(net.ParseIP("192.168.66.6"))
func IsPublicIP(IP net.IP) bool {
	if IP.IsLoopback() || IP.IsLinkLocalMulticast() || IP.IsLinkLocalUnicast() {
		return false
	}
	if ip4 := IP.To4(); ip4 != nil {
		switch {
		case ip4[0] == 10:
			return false
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return false
		case ip4[0] == 192 && ip4[1] == 168:
			return false
		default:
			return true
		}
	}
	return false
}
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
	firewall, err := FirewallAllowMyCurrentIp(ctx, client, firewallName, tagName)
	if err != nil {
		// Output Firewall information
		logrus.WithFields(logrus.Fields{
			"ID": firewall.ID,
			// "Name":          firewall.Name,
			// "Status":        firewall.Status,
			"InboundRules":  firewall.InboundRules,
			"OutboundRules": firewall.OutboundRules,
		}).Error(fmt.Sprintf("firewall '%s' update %s", firewall.Name, firewall.Status))
	}

	// Output Firewall information
	logrus.WithFields(logrus.Fields{
		"ID": firewall.ID,
		// "Name":          firewall.Name,
		// "Status":        firewall.Status,
		// "InboundRules":  firewall.InboundRules,
		// "OutboundRules": firewall.OutboundRules,
	}).Info(fmt.Sprintf("firewall '%s' update %s", firewall.Name, firewall.Status))
}

func main() {
	// Run command $(make load-envs)
	stdout_bytes, err := exec.Command("make", "-s", "load-envs").Output()
	if err != nil {
		log.Error(err)
	}
	stdout := string(stdout_bytes)[:len(stdout_bytes)-1]
	log.Debug("Envs from makefile: " + stdout)

	// Parse envs to a list from $(make load-envs)
	for _, env := range strings.Split(stdout, " ") {
		env_splits := strings.Split(env, "=")
		if len(env_splits) > 1 && env_splits[1] != "" {
			os.Setenv(strings.TrimSpace(env_splits[0]), strings.TrimSpace(env_splits[1]))
		}
	}

	// Check if env var SENTRY_DSN is set
	if os.Getenv("SENTRY_DSN") != "" {
		// Initialize Sentry
		sentry_dsn := os.Getenv("SENTRY_DSN")
		logger, err := log.NewSession(sentry_dsn)

		if err != nil {
			fmt.Println("[ERROR] Can't initialize Sentry Logger!", err)
		}

		// Set panic handler
		defer func() {
			if err := recover(); err != nil {
				logger.Panic(err)
			}
		}()
	}

	// log.Panic("test")

	// os.Exit(0)

	// Run updateFirewall once
	updateFirewall()

	// Add updateFirewall to scheduler
	gocron.Every(30).Minutes().Do(updateFirewall)

	// Start scheduler
	<-gocron.Start()
}
