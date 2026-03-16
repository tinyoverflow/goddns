package main

import (
	"flag"
	"fmt"
	"goddns/internal/provider"
	"goddns/internal/retriever"
	"log"
)

func main() {
	ret := retriever.IfConfigRetriever{}
	ip, err := ret.GetIPAddress()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ip)

	var hetznerToken string
	flag.StringVar(&hetznerToken, "token", "", "Hetzner Cloud Token")

	var dnsZone string
	flag.StringVar(&dnsZone, "zone", "", "DNS Zone")

	p := provider.HetznerCloudProvider{
		Token: hetznerToken,
		Zone:  dnsZone,
	}

	if p.SetIPAddress(ip) != nil {
		log.Fatal(err)
	}
}
