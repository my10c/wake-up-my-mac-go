//
// BSD 3-Clause License
//
// Copyright © 2023, © Badassops LLC / Luc Suryo
// All rights reserved.
//

package validator

import (
	"fmt"
	"net"
	"os"

	"configurator"

	"github.com/my10c/packages-go/print"

	valid "github.com/muonsoft/validation/validate"
)

func ValidateConfig(c *configurator.Config, p *print.Print) {
	var errMsg string
	var ipOK bool = false
	var ips []string
	var err error

	// MAC address check, must always be set
	if _, err = net.ParseMAC(c.MACAddress); err != nil {
		errMsg = fmt.Sprintf("Invalid MAC addess %s\n", c.MACAddress)
		p.PrintRed(errMsg)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// If IP was set, we need to make sure is either a valid IPv4 or IPv6 address
	if len(c.IPAddress) > 0 {
		// check IPv4 first
		if valid.IPv4(c.IPAddress) == nil {
			if c.Debug {
				p.PrintYellow(fmt.Sprintf("\tUsing resolved IPv4 %v\n", c.IPAddress))
			}
			c.IPv = 4
			ipOK = true
		}
		// previoud check was not a valid IPv4 now check for valid IPv4
		if !ipOK {
			// check IPv6
			if valid.IPv6(c.IPAddress) != nil {
				// IP address was not valid
				errMsg = fmt.Sprintf("IP %s is not a valid IPv6 nor a IPv4 address\n", c.IPAddress)
				p.PrintRed(errMsg)
				os.Exit(1)
			}
			if c.Debug {
				p.PrintYellow(fmt.Sprintf("\tUsing resolved  %v\n", c.IPAddress))
			}
			c.IPv = 6
		}
	}

	// in case no ip was give we need to be able to get IP from hostname
	// 1. resolved the given host name
	// 2. use the resolver to get the IPv4 or IPv6 addresses
	if len(c.HostName) > 0 && len(c.IPAddress) == 0 {
		if ips, err = net.LookupHost(c.HostName); err != nil {
			errMsg = fmt.Sprintf("Unable to resolve host %s\n", c.HostName)
			p.PrintRed(errMsg)
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		ipOK = false
		// check all IP addresses, and stop at the first valid IP (IPv4 or IPv6)
		for idx := 0; idx < len(ips); idx++ { 
			// we prefer IPv4
			if valid.IPv4(ips[idx]) == nil {
				if c.Debug {
					p.PrintYellow(fmt.Sprintf("\tUsing resolved IPv4 %v\n", ips[idx]))
				}
				c.IPAddress = ips[idx]
				c.IPv = 4
				ipOK = true
				break
			}
			if valid.IPv6(ips[idx]) == nil {
				if c.Debug {
					p.PrintYellow(fmt.Sprintf("\tUsing resolved IPv6 %v\n", ips[idx]))
				}
				c.IPAddress = ips[idx]
				c.IPv = 6
				ipOK = true
				break
			}
		}
		// unable to get IP from hostname
		if !ipOK {
			errMsg = fmt.Sprintf("Unable to resolve hostname %v to an IPv4 or IPv6 address\n", ips)
			p.PrintRed(errMsg)
			os.Exit(1)
		}
	}

	// if host was not set. let tro reverse lookup
	if len(c.HostName) == 0 {
		ips, err = net.LookupAddr(c.IPAddress)
		if err != nil {
			errMsg = fmt.Sprintf("Unable to resolve ip %s to host\n", c.IPAddress)
            p.PrintYellow(errMsg)
		}
		// we use first entry
		c.HostName = ips[0]
	}
}
