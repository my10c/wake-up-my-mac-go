//
// BSD 3-Clause License
//
// Copyright © 2023, © Badassops LLC / Luc Suryo
// All rights reserved.
//

package vars

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"time"
)

var (
	Off		= "\x1b[0m"		// Text Reset
	Black	= "\x1b[1;30m"	// Black
	Red		= "\x1b[1;31m"	// Red
	Green	= "\x1b[1;32m"	// Green
	Yellow	= "\x1b[1;33m"	// Yellow
	Blue	= "\x1b[1;34m"	// Blue
	Purple	= "\x1b[1;35m"	// Purple
	Cyan	= "\x1b[1;36m"	// Cyan
	White	= "\x1b[1;37m"	// White

	RedUnderline	= "\x1b[4;31m" // Red underline
	OneLineUP		= "\x1b[A"
)

var (
	MyVersion	= "0.0.1"
	now			= time.Now()
	MyProgname	= path.Base(os.Args[0])
	myAuthor	= "Luc Suryo"
	myCopyright = "Copyright 2023 - " + strconv.Itoa(now.Year()) + " ©Badassops LLC"
	myLicense	= "License 3-Clause BSD, https://opensource.org/licenses/BSD-3-Clause ♥"
	myEmail		= "<luc@badassops.com>"
	MyInfo		= fmt.Sprintf("%s (version %s)\n%s\n%s\nWritten by %s %s\n",
		MyProgname, MyVersion, myCopyright, myLicense, myAuthor, myEmail)
	MyDescription = "Script to send a wake up signal to a computer 😁"

	// default wait 
	Wait = 20

	// default port
	Port = 80

	// default config file
	WakeFile = "/usr/local/etc/wake/wake.conf"
)
