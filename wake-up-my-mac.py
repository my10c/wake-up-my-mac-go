#!/usr/bin/env python3
# coding: utf8

""" A python3 script to wake up a Apple computer """

# BSD 3-Clause License
#
# Copyright (c) 2018 - 2023, © Badassops LLC / Luc Suryo
# All rights reserved.
#
# Redistribution and use in source and binary forms, with or without
# modification, are permitted provided that the following conditions are met:
#
# * Redistributions of source code must retain the above copyright notice, this
#   list of conditions and the following disclaimer.
#
# * Redistributions in binary form must reproduce the above copyright notice,
#   this list of conditions and the following disclaimer in the documentation
#   and/or other materials provided with the distribution.
#
# * Neither the name of the copyright holder nor the names of its
#   contributors may be used to endorse or promote products derived from
#   this software without specific prior written permission.
#
# THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
# AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
# IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
# DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
# FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
# DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
# SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
# CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
# OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
# OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
#*
#* File        :   wake_mac
#*
#* Description :   script to send wake signal up to an Apple  computer
#*
#* Author      :    Luc Suryo <luc@badassops.com>
#*
#* Version     :    0.2
#*
#* Date        :    Aug 1, 2023
#*
#* History    :
#*    Date:         Author:    Info:
#*    Dec 1, 2018   LIS        First Release (private)
#*    Aug 1, 2023   LIS        Public release without hardcode values
#
# TODO : make it work and get the -c option coded

# for debugging
# import pdb
# pdb.set_trace()

import socket
import signal
import sys
import threading
import unicodedata
import os
import argparse
import ipaddress
import macaddress
import re
from time import time, sleep, strftime

# COLORS
OFF='\033[0m'        # Text Reset

# Bold
BLACK  = '\033[1;30m'
RED    = '\033[1;31m'
GREEN  = '\033[1;32m'
YELLOW = '\033[1;33m'
BLUE   = '\033[1;34m'
PURPLE = '\033[1;35m'
CYAN   = '\033[1;36m'
WHITE  = '\033[1;37m'

__progname__ = os.path.basename(__file__)
__author__ = 'Luc Suryo'
__copyright__ = 'Copyright 2019 - ' + strftime('%Y') + ' © Badassops LLC'
__license__ = 'License 3-Clause BSD, https://opensource.org/licenses/BSD-3-Clause ♥'
__version__ = '0.2'
__email__ = '<luc@badassops.com>'
__info__ = '%s%s\n%s\nLicense %s\n\nWritten by %s %s%s\n' % \
        (YELLOW, __version__, __copyright__, __license__, __author__, __email__, OFF)
### TODO __usage_txt__  = '%s<[--mac=] | [--hostname=]> [--ip=] | [--config=]%s' %\
__usage_txt__  = '%s[-h] [-v] -m MAC (-i IP | -H HOSTNAME) [-w WAIT]%s' % (PURPLE, OFF)
__description__ = '%sScript to wake up an Apple computer%s' % (YELLOW, OFF)

# DEFAULTS
SLEEP = 20

# Help messages
MAC_HELP = '%sMAC address of the Apple computer (required)%s' % (BLUE, OFF)
IP_HELP = '%sValid IPv4 address, can not be used with -H%s' % (BLUE, OFF)
HOSTNAME_HELP = '%sDNS hostname of the Apple computer, can not be used with -i%s'  % (BLUE, OFF)
### TODO CONFIG_HELP = '%sConfiguration file with 2 entries, 1 per lines;  MAC and IP or HOST%s'
WAIT_HELP = '%sHow many seconds to wait, default to %s%s' % (BLUE, SLEEP, OFF)

# we use port 80 to wake up the computer
PORT = 80

def signal_handler(signum, frame):
    """Signal/interrupts handler
        @param  signum  {int}       The interrupt ID according to signal.h.
        @param  frame   {string}    Memory frame where the interrupted was called.
    """

    if signum is signal.SIGHUP:
        print ('\n\tProcess kill -HUP received.')
    elif signum is signal.SIGINT :
        print ('\n\tProcess aborted on your request, ctrl-c received.')
    elif signum is signal.SIGTERM :
        print ('\n\tProcess kill -TERM received.')
    else:
        print ('\n\tProcess aborted due to received signal: %d.' % signum)

    SPIN.stop()
    sys.exit(128)

class SpinCursor(threading.Thread):
    """ Class and function so display a wait spinner (dots or wheel)
    """

    def __init__(self, msg=None, counter=0, maxspin=0, minspin=10, speed=5, dots=False):
        # Count of a spin
        self.count = 0
        self.out = sys.stdout
        self.flag = False
        self.max = maxspin
        self.min = minspin

        # Any message to print first ?
        self.msg = msg

        # Complete printed string
        self.string = None

        # counter
        self.counter = counter

        # Speed is given as number of spins a second
        # Use it to calculate spin wait time
        self.waittime = 1.0/float(speed*4)
        if os.name == 'posix':
            if dots is True:
                self.spinchars = (unicodedata.lookup('FIGURE DASH'), u'. ', u'• ', u'. ', u'• ')
            else:
                self.spinchars = (unicodedata.lookup('FIGURE DASH'), u'\\ ', u'| ', u'/ ')
        else:
            # The unicode dash character does not show
            # up properly in Windows console.
            if dots is True:
                self.spinchars = (u'. ', u'• ', u'. ', u'• ')
            else:
                self.spinchars = (u'-', u'\\ ', u'| ', u'/ ')
        threading.Thread.__init__(self, None, None, "Spin Thread")

    def spin(self):
        """ Perform a single spin """
        for spinchar in self.spinchars:
            if self.msg:
                self.string = self.msg + '...\t' + spinchar + '\r'
            else:
                self.string = '...\t' + spinchar + '\r'
            self.out.write(self.string)
            self.out.flush()
            sleep(self.waittime)

    def run(self):
        """ run spinning """
        while (not self.flag) and ((self.count < self.min) or (self.count < self.max)):
            self.spin()
            self.count += 1
        # Clean up display...
        self.out.write('\033[2K')

    def spinCounter(self):
        """ Perform a single spin """
        print('\033[2K {}: [{}]'.format(self.msg, str(self.counter)), end='\r', flush=True)
        sleep(self.waittime * 5)

    def runCounter(self):
        """ run spinning """
        print('{}'.format(PURPLE), end='\r', flush=True)
        while (not self.flag) and (self.counter  != 0):
            self.spinCounter()
            self.counter -= 1
        # Clean up display...
        self.out.write('\033[2K')

    def stop(self):
        """ stop spinning """
        print('{}'.format(OFF), end='\r', flush=True)
        self.flag = True

def checkHost(hostname):
    try:
        socket.gethostbyname(hostname)
        return True
    except socket.error:
        return False

if __name__ == "__main__":
    argOK = True

    # Install signal/interrupts handler, we capture only SIGHUP, SIGINT and TERM
    signal.signal(signal.SIGHUP, signal_handler)
    signal.signal(signal.SIGINT, signal_handler)
    signal.signal(signal.SIGTERM, signal_handler)

    # Process giving arguments
    PARSER = argparse.ArgumentParser(usage=__usage_txt__, description=__description__, \
        formatter_class=argparse.RawDescriptionHelpFormatter, conflict_handler='resolve')

    IPHOSTNAME = PARSER.add_mutually_exclusive_group(required=True)

    PARSER.add_argument('-v', '--version',
        action='version', version=__info__)

    PARSER.add_argument('-m', '--mac',
        action='store', dest='mac', required=True, help=MAC_HELP)

    IPHOSTNAME.add_argument('-i', '--ip',
        action='store', dest='ip', help=IP_HELP)

    IPHOSTNAME.add_argument('-H', '--hostname',
        action='store', dest='hostname', help=HOSTNAME_HELP)

    PARSER.add_argument('-w', '--wait',
        action='store', dest='wait', required=False, type=int, default=SLEEP, help=WAIT_HELP)

    ARGS = PARSER.parse_args()

    if ARGS.ip is not None:
        try:
            ipaddress.IPv4Address(ARGS.ip)
            tulpValue = ARGS.ip
        except ValueError:
            print('Given IP address is invalid: %s' % ARGS.ip)
            argOK = False

    if ARGS.hostname is not None:
        tulpValue = ARGS.hostname
        if checkHost(ARGS.hostname) == False:
            print('Given hostname is not resolvable  %s' % ARGS.hostname)
            argOK = False

    mac = ARGS.mac.lower()
    if not bool(re.match('^' + '[\:\-]'.join(['([0-9a-f]{2})']*6) + '$', mac)):
        print('Given MAC address is invalid: %s' % ARGS.mac)
        argOK = False

    if not argOK:
        sys.exit(1)

    # configure the socket
    SOCKET = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    SOCKET.setsockopt(socket.SOL_SOCKET, socket.SO_BROADCAST, 1)

    # convert given MAC-address into bytes
    macBytes = bytes(macaddress.MAC(mac))

    ### TODO make this work?

    try:
        SOCKET.sendto(b'\xff'*6 + macBytes  *16, (tulpValue, PORT))
    except Exception as err:
        print('Error sending the wake up signal, error {}'.format(err))
        sys.exit(1)

    SPIN = SpinCursor(msg='sleeping for {} seconds'.format(ARGS.wait), counter=ARGS.wait, speed=1)
    SPIN.runCounter()
    SPIN.stop()
    print('{}The Apple computer ({}) with MAC address {} should be awake now.{}\n'
            .format(GREEN, tulpValue, ARGS.mac, OFF))
    sys.exit(0)
