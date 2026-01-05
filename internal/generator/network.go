/*-------------------------------------------------------------------------
 *
 * pgEdge Anonymizer
 *
 * Portions copyright (c) 2025 - 2026, pgEdge, Inc.
 * This software is released under The PostgreSQL License
 *
 *-------------------------------------------------------------------------
 */

package generator

import (
	"fmt"
	"strings"

	"github.com/pgedge/pgedge-anonymizer/internal/generator/data"
)

// IPv4Generator generates IPv4 addresses.
type IPv4Generator struct {
	BaseGenerator
}

// NewIPv4Generator creates a new IPv4 address generator.
func NewIPv4Generator() *IPv4Generator {
	return &IPv4Generator{
		BaseGenerator: BaseGenerator{name: "IPV4_ADDRESS"},
	}
}

// Generate produces an IPv4 address.
// It avoids reserved ranges and generates realistic-looking addresses.
func (g *IPv4Generator) Generate(input string) string {
	// Generate random octets, avoiding reserved ranges
	// Use common private ranges or realistic public-looking addresses
	firstOctet := g.randomFirstOctet()
	return fmt.Sprintf("%d.%d.%d.%d",
		firstOctet,
		randomInt(256),
		randomInt(256),
		1+randomInt(254), // Avoid .0 and .255
	)
}

// randomFirstOctet generates a valid first octet, avoiding problematic ranges.
func (g *IPv4Generator) randomFirstOctet() int {
	// Choose between private ranges and realistic public ranges
	choice := randomInt(5)
	switch choice {
	case 0:
		return 10 // 10.x.x.x (private)
	case 1:
		return 172 // Could be 172.16-31.x.x (private)
	case 2:
		return 192 // Could be 192.168.x.x (private)
	default:
		// Generate a "public-looking" first octet
		// Avoid 0, 127 (loopback), 224-255 (multicast/reserved)
		for {
			octet := 1 + randomInt(223) // 1-223
			if octet != 127 && octet != 10 {
				return octet
			}
		}
	}
}

// IPv6Generator generates IPv6 addresses.
type IPv6Generator struct {
	BaseGenerator
}

// NewIPv6Generator creates a new IPv6 address generator.
func NewIPv6Generator() *IPv6Generator {
	return &IPv6Generator{
		BaseGenerator: BaseGenerator{name: "IPV6_ADDRESS"},
	}
}

// Generate produces an IPv6 address.
// It detects the input format and generates a matching format.
func (g *IPv6Generator) Generate(input string) string {
	// Detect if input uses compressed format (::)
	compressed := strings.Contains(input, "::")

	// Detect if input uses uppercase
	uppercase := strings.ToUpper(input) == input && strings.ContainsAny(input, "ABCDEF")

	// Generate 8 groups of 4 hex digits
	groups := make([]string, 8)
	for i := 0; i < 8; i++ {
		groups[i] = g.randomHexGroup(uppercase)
	}

	if compressed {
		// Use compressed format - generate with :: notation
		if randomInt(2) == 0 {
			return fmt.Sprintf("2001:db8:%s:%s::%s",
				g.randomHexGroup(uppercase),
				g.randomHexGroup(uppercase),
				g.randomHexGroup(uppercase))
		}
	}

	return strings.Join(groups, ":")
}

// randomHexGroup generates a random 4-character hex group.
func (g *IPv6Generator) randomHexGroup(uppercase bool) string {
	chars := "0123456789abcdef"
	if uppercase {
		chars = "0123456789ABCDEF"
	}

	result := make([]byte, 4)
	for i := 0; i < 4; i++ {
		result[i] = chars[randomInt(16)]
	}
	return string(result)
}

// HostnameGenerator generates hostnames.
type HostnameGenerator struct {
	BaseGenerator
	data *data.DataSet
}

// NewHostnameGenerator creates a new hostname generator.
func NewHostnameGenerator(d *data.DataSet) *HostnameGenerator {
	return &HostnameGenerator{
		BaseGenerator: BaseGenerator{name: "HOSTNAME"},
		data:          d,
	}
}

// hostname prefixes for generating realistic hostnames
var hostnamePrefixes = []string{
	"server", "srv", "web", "www", "app", "api", "db", "mail", "mx",
	"ns", "dns", "ftp", "vpn", "gateway", "gw", "proxy", "cache",
	"node", "worker", "master", "slave", "primary", "replica",
	"dev", "staging", "prod", "test", "qa", "uat",
	"host", "vm", "container", "k8s", "docker",
	"linux", "win", "ubuntu", "centos", "debian",
	"us-east", "us-west", "eu-west", "ap-south",
}

// hostname domains for generating realistic FQDNs
var hostnameDomains = []string{
	"example.com", "example.org", "example.net",
	"internal", "local", "localdomain", "corp", "private",
	"cloud.local", "datacenter.local", "cluster.local",
}

// Generate produces a hostname.
// It detects the input format and generates a matching style.
func (g *HostnameGenerator) Generate(input string) string {
	// Check if input is a FQDN (contains dots)
	isFQDN := strings.Contains(input, ".")

	// Check if input has numeric suffix
	hasNumber := false
	for _, c := range input {
		if c >= '0' && c <= '9' {
			hasNumber = true
			break
		}
	}

	// Generate hostname
	prefix := hostnamePrefixes[randomInt(len(hostnamePrefixes))]

	var hostname string
	if hasNumber {
		// Add numeric suffix
		hostname = fmt.Sprintf("%s%02d", prefix, 1+randomInt(99))
	} else {
		hostname = prefix
	}

	if isFQDN {
		// Add domain
		domain := hostnameDomains[randomInt(len(hostnameDomains))]
		return hostname + "." + domain
	}

	return hostname
}
