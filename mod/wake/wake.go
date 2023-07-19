//
// BSD 3-Clause License
//
// Copyright © 2023, © Badassops LLC / Luc Suryo
// All rights reserved.
//

package wake

import (
	"bytes"
    "encoding/binary"
	"fmt"
	"net"
	"os"
	"syscall"

	"configurator"
	"github.com/my10c/packages-go/print"
)

type (
 	MACAddress [6]byte
 
 	WakeUpPacket struct {
 		header [6]byte
 		payload [16]MACAddress
 	}
)

var (
	err error
)

func setupNetwork(c *configurator.Config, p *print.Print) (int) {
	var fd int
	// setup the connection
	if c.IPv == 4 {
		fd, err = syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	}

	if c.IPv == 6 {
		fd, err = syscall.Socket(syscall.AF_INET6, syscall.SOCK_DGRAM, 0)
	}

	if err != nil {
		p.PrintRed("Error creating socket")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_BROADCAST, 1)
	if err != nil {
		p.PrintRed("Error setting the socket options")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return fd
}

func setupPayLoad(c *configurator.Config, p *print.Print) bytes.Buffer {
	// now ready to starts step the wakeup packet
 	var macBytes [6]byte
 	var packet WakeUpPacket
	var payload bytes.Buffer

	// Setup the header which is 6 repetitions of 0xFF.
	for idx := range packet.header {
		packet.header[idx] = 0xFF
	}

	// Convert the MAC address into byte
	for idx := range macBytes {
		macBytes[idx] = c.MACAddress[idx]
	}

	// Setup the payload which is 16 repetitions of the MAC addr
	for idx := range packet.payload {
		packet.payload[idx] = macBytes
	}

	// create the payload
	if err := binary.Write(&payload, binary.BigEndian, packet); err != nil {
		p.PrintRed("Error setting the socket options")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return payload
}

func WakeMe(c *configurator.Config, p *print.Print) {
	// create the socket address
	fd :=  setupNetwork(c, p)

	// create the payload
	payload :=  setupPayLoad(c, p)

	p.PrintGreen(fmt.Sprintf("\tReady to send the wake-up pakket IPv%d\n", c.IPv))
	if c.IPv == 6 {
		sockaddr := new(syscall.SockaddrInet6)
		for i := 0; i < net.IPv6len; i++ {
			sockaddr.Addr[i] = c.IPAddress[i]
		}
		// ZoneId uint32
		sockaddr.Port = c.Port
		err = syscall.Sendto(fd, payload.Bytes(), 0, sockaddr)
	}
	if c.IPv == 4 {
		sockaddr := new(syscall.SockaddrInet4)
		for i := 0; i < net.IPv4len; i++ {
			sockaddr.Addr[i] = c.IPAddress[i]
		}
		sockaddr.Port = c.Port
		err = syscall.Sendto(fd, payload.Bytes(), 0, sockaddr)
	}

	if err != nil {
		p.PrintRed("Error sending the wake up packet")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
