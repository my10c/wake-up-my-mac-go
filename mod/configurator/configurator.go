//
// BSD 3-Clause License
//
// Copyright © 2023, © Badassops LLC / Luc Suryo
// All rights reserved.
//

package configurator

import (
	"fmt"
	"os"

	"vars"

	"github.com/my10c/packages-go/is"
	"github.com/my10c/packages-go/print"

	"github.com/akamensky/argparse"
	"github.com/BurntSushi/toml"
)

type (
	Config struct {
		ConfigFile	string
		MACAddress	string
		IPAddress	string
		HostName	string
		Port		int
		Wait		int
		Debug		bool
		IPv			int
	}

	Host struct {
		HostName	string	`toml:"hostname,omitempty"`
		MACAddress	string	`toml:"mac,omitempty"`
		IPAddress	string	`toml:"ip,omitempty"`
		Port		int		`toml:"port,omitempty"`
		Wait		int		`toml:"wait,omitempty"`
	}

	tomlConfig struct {
		Hosts		map[string]Host	`toml:"hosts,omitempty"`
	}
)

var (
	defaultSet		bool
	configFileSet	bool
	macAddressSet	bool
	ipAddressSet	bool
	hostNameSet		bool
	portSet			bool
	waitSet			bool
)

// function to initialize the configuration
func Configurator() *Config {
	// the rest of the values will be filled from the given configuration file
	return &Config{
		Debug: false,
		IPv: 4,
	}
}

func (c *Config) InitializeArgs(p *print.Print, i *is.Is) {
	parser := argparse.NewParser(vars.MyProgname, vars.MyDescription)

	wakeFile := parser.Flag("C", "default",
		&argparse.Options{
			Required:	false,
			Help:		fmt.Sprintf("Use the default config file %s", vars.WakeFile),
			Default:	false,
		})

	configFile := parser.String("c", "configFile",
		&argparse.Options{
			Required: false,
			Help:		"Configuration file to be use",
		})

	macAddress := parser.String("m", "macAddress",
		&argparse.Options{
			Required: false,
			Help:		"The MAC address of the computer, can not be use with the -c flag",
		})

	ipAddress := parser.String("i", "ipAddress",
		&argparse.Options{
			Required:	false,
			Help:		"The IP address of the computer, can not be use with the -c flag",
		})

	hostName := parser.String("H", "host",
		&argparse.Options{
			Required:	false,
			Help:		"The hostname of the computer, required with the -c flag",
		})

	port := parser.Int("p", "port",
		&argparse.Options{
			Required:	false,
			Help:		"Computer network port to send the wake up signal",
			Default:	vars.Port,
		})

	wait := parser.Int("w", "wait",
		&argparse.Options{
			Required:	false,
			Help:		"How many seconds to wait for the computer to be awake",
			Default:	vars.Wait,
		})

	debug := parser.Flag("d", "debug",
		&argparse.Options{
			Required:	false,
			Help:		"Enable some debug output",
			Default:	false,
		})

	showVersion := parser.Flag("v", "version",
		&argparse.Options{
		Required:	false,
		Help:		"Show version",
	})

	err := parser.Parse(os.Args)
	if err != nil {
		p.PrintBlue(parser.Usage(err))
		os.Exit(1)
	}

	if *showVersion {
		p.ClearScreen()
		p.PrintYellow(vars.MyProgname + " version: " + vars.MyVersion + "\n")
		os.Exit(0)
	}

	if *debug {
		c.Debug = true
	}

	// get the flags that were set
	defaultSet = parser.GetArgs()[1].GetParsed()
	configFileSet = parser.GetArgs()[2].GetParsed()
	macAddressSet = parser.GetArgs()[3].GetParsed()
	ipAddressSet = parser.GetArgs()[4].GetParsed()
	hostNameSet = parser.GetArgs()[5].GetParsed()
	portSet = parser.GetArgs()[6].GetParsed()
	waitSet = parser.GetArgs()[7].GetParsed()

	// handle wrong combination
	// -C and -c can not be used together
	if defaultSet && configFileSet {
		p.PrintBlue(vars.MyInfo)
		p.PrintRed("\n\tThe flags -c and -C are mutually exclusive\n\n")
		p.PrintGreen(parser.Usage(err))
		os.Exit(1)
	}

	// -C requires only -H and no other flags
	if (defaultSet) {
		if (!hostNameSet) || (macAddressSet || ipAddressSet) {
			p.PrintBlue(vars.MyInfo)
			p.PrintRed("\n\tThe flag -C requires only the -H flags, but not -m or -i\n\n")
			p.PrintGreen(parser.Usage(err))
			os.Exit(1)
		}
	}

	// -c requires only -H can not be combined with -m or ip
	if (configFileSet) {
		if (!hostNameSet) || (macAddressSet || ipAddressSet) {
			p.PrintBlue(vars.MyInfo)
			p.PrintRed("\n\tThe flag -c requires only the -H flags, but not -m or -i\n\n")
			p.PrintGreen(parser.Usage(err))
			os.Exit(1)
		}
	}

	// if -C or -c was not given it requires the -m and at least -i or -H or both
	if (!defaultSet && !configFileSet) {
		if !macAddressSet || (!ipAddressSet && !hostNameSet) {
			p.PrintBlue(vars.MyInfo)
			p.PrintRed("\n\tWhen not using the -C or -C the flags,\n")
			p.PrintRed("\tflags -m is required with either the -i or -H flags, or both\n\n")
			p.PrintGreen(parser.Usage(err))
			os.Exit(1)
		}
	}

	// get given values
	c.ConfigFile = *configFile
	// if -C was given we overwrite the config file to use
	if *wakeFile {
		c.ConfigFile = vars.WakeFile
	}
	// make sure the given file or the default exist if (-C and -c flag)
	if len(c.ConfigFile) > 0 {
		if _, ok, _ := i.IsExist(c.ConfigFile, "file"); !ok {
			p.PrintBlue(vars.MyInfo)
			p.PrintRed("\nThe given configuration file " + c.ConfigFile + " does not exist\n\n")
			os.Exit(1)
		}
	}

	c.MACAddress = *macAddress
	c.IPAddress = *ipAddress
	c.HostName = *hostName
	c.Port = *port
	c.Wait = *wait
	// if -C or -c was given we need to set values from file
	if len(c.ConfigFile) > 0 {
		c.getConfigFromFile(p)
	}
}

func (c *Config) getConfigFromFile(p *print.Print) {
	var configValues tomlConfig

	if _, err := toml.DecodeFile(c.ConfigFile, &configValues); err != nil {
		p.PrintRed("Error reading the configuration file\n")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// we always need the MAC address
	c.MACAddress = configValues.Hosts[c.HostName].MACAddress
	if len(c.MACAddress) == 0 {
		errMsg := fmt.Sprintf("make sure the host %s has an entry in the config file %s\n",
			c.HostName, c.ConfigFile,
		)
		p.PrintRed("Configuration file error, MAC entry is missing\n")
		p.PrintRed(errMsg)
		os.Exit(1)
	}

	// the IP address
	c.IPAddress = configValues.Hosts[c.HostName].IPAddress
	if len(c.IPAddress) == 0 {
		p.PrintRed("Configuration file error, IP entry is missing\n")
		os.Exit(1)
	}

	// Port and wait are not use if it was given on the cli
	if !portSet {
		if configValues.Hosts[c.HostName].Port != 0 {
			c.Port = configValues.Hosts[c.HostName].Port
		}
	}

	if !waitSet {
		if configValues.Hosts[c.HostName].Wait != 0 {
			c.Wait = configValues.Hosts[c.HostName].Wait
		}
	}

	// the hostname always as last, if we overwrite it, we can no longer
	// use to scan the value from the config file; [c.HostName]
	if len(configValues.Hosts[c.HostName].HostName) > 0 {
		c.HostName = configValues.Hosts[c.HostName].HostName
	}
}
