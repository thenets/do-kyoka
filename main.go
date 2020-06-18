package main

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/spf13/viper"
	"github.com/thenets/do-kyoka/helper"
)

func main() {
	// Read config file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/do-kyoka/")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	// Params
	firewallName := viper.GetString("firewall.name")
	if firewallName == "" {
		panic("[ERROR] 'firewall.name' not set in config file!")
	}
	tagName := viper.GetString("firewall.tag")
	if tagName == "" {
		panic("[ERROR] 'firewall.tag' not set in config file!")
	}
	apiToken := viper.GetString("apiToken")
	if tagName == "" {
		panic("[ERROR] 'apiToken' not set in config file!")
	}

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
