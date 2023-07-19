# Under construction

# wake-up-my-mac
Script to wake up a Apple computer via network

## Python version
```
usage: wake-up-my-mac [-h|--help] [-c|--configFile "<value>"] [-m|--macAddress
                      "<value>"] [-i|--ipAddress "<value>"] [-H|--host
                      "<value>"] [-C|--default] [-p|--port <integer>]
                      [-w|--wait <integer>] [-d|--debug] [-v|--version]

                      Script to send a wake up signal to a computer 😁

Arguments:

  -h  --help        Print help information
  -c  --configFile  Configuration file to be use
  -m  --macAddress  The MAC address of the computer, can not be use with the -c
                    flag
  -i  --ipAddress   The IP address of the computer, can not be use with the -c
                    flag
  -H  --host        The hostname of the computer, required with the -c flag
  -C  --default     Use the default config file /usr/local/etc/wake/wake.conf.
                    Default: false
  -p  --port        Computer network port to send the wake up signal. Default:
                    80
  -w  --wait        How many seconds to wait for the computer to be awake.
                    Default: 20
  -d  --debug       Enable some debug output. Default: false
  -v  --version     Show version

```

### Configuration file
- only for the Go version
- see an example under the example directory

### Notes
- I use this with my Apple computers, but not sure is common enough that it will work with any other computer

### TODO
- Support of IPv6 not tested
