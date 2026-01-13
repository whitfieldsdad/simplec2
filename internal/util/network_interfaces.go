package util

import (
	"context"
	"net"
)

type NetworkInterface struct {
	Name          string   `json:"name"`
	IPv4Addresses []string `json:"ipv4_addresses"`
	IPv6Addresses []string `json:"ipv6_addresses"`
	MACAddress    string   `json:"mac_address"`
}

func (nic NetworkInterface) GetArtifactType() ArtifactType {
	return ArtifactTypeNetworkInterface
}

func ListNetworkInterfaces(ctx context.Context) ([]NetworkInterface, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	nics := []NetworkInterface{}
	for _, iface := range ifaces {
		ipv4Addresses := []string{}
		ipv6Addresses := []string{}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			ip, _, err := net.ParseCIDR(addr.String())
			if err != nil {
				continue
			}
			if ip.To4() != nil {
				ipv4Addresses = append(ipv4Addresses, ip.String())
			}
			if ip.To16() != nil && ip.To4() == nil {
				ipv6Addresses = append(ipv6Addresses, ip.String())
			}
		}
		nic := NetworkInterface{
			Name:          iface.Name,
			IPv4Addresses: ipv4Addresses,
			IPv6Addresses: ipv6Addresses,
			MACAddress:    iface.HardwareAddr.String(),
		}
		nics = append(nics, nic)
	}
	return nics, nil
}
