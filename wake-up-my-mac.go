// BSD 3-Clause License
//
// Copyright © 2023, © Badassops LLC / Luc Suryo
// All rights reserved.
//
// Version	:	0.1
//

package main

import (
	"fmt"
	"os"
	"strconv"

	// local
	"configurator"
	"validator"
	"wake"

	// on github
	"github.com/my10c/packages-go/is"
	"github.com/my10c/packages-go/print"
	"github.com/my10c/packages-go/spinner"
)

var (
	configData []string
)

func main() {
	var endMesg string
	s := spinner.New(1000)
	p := print.New()
	i := is.New()
	c := configurator.Configurator()

	// get given parameters
	c.InitializeArgs(p, i)

	fmt.Printf("\t%s\n", p.PrintLine(print.Purple, 60))
	// check for valid parameters
	validator.ValidateConfig(c, p)

	if c.Debug {
		p.PrintBlue(fmt.Sprintf("\tIP         : %v\n", c.IPAddress))
		p.PrintBlue(fmt.Sprintf("\tHostname   : %v\n", c.HostName))
		p.PrintBlue(fmt.Sprintf("\tMAC        : %v\n", c.MACAddress))
		if len(c.ConfigFile) > 0 {
			p.PrintBlue(fmt.Sprintf("\tConfigfile : %v\n", c.ConfigFile))
		}
	}
	// send the wake singal
	wake.WakeMe(c, p)

	// counter count down
	 s.Counter(c.Wait, "\tWaiting for " + strconv.FormatInt(int64(c.Wait), 10) + " seconds for the compter to wake up")

	// add Hostname is it was set
	if len(c.HostName) != 0 {
		endMesg = fmt.Sprintf("\tThe computer %s, IP: %s\n\tMAC address: %s should be awake now\n",
			c.HostName, c.IPAddress, c.MACAddress,
		)
	} else {
		endMesg = fmt.Sprintf("\tThe computer, IP: %s\n\tMAC address: %s should be awake now\n",
			c.HostName, c.MACAddress,
		)
	}

	fmt.Printf("\t%s\n", p.PrintLine(print.Purple, 60))
	p.PrintGreen(endMesg)
	fmt.Printf("\t%s\n", p.PrintLine(print.Purple, 60))
	p.TheEnd()
	fmt.Printf("\t%s\n", p.PrintLine(print.Purple, 60))

	os.Exit(0)
}
