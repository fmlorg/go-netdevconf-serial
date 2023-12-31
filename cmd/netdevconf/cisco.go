// -*- utf-8 -*-
//
// Copyright (C) 2023 Ken'ichi Fukamachi
//   All rights reserved. This program is free software; you can
//   redistribute it and/or modify it under 2-Clause BSD License.
//   https://opensource.org/licenses/BSD-2-Clause
//
// mailto: fukachan@fml.org
//    web: https://www.fml.org/
// github: https://github.com/fmlorg
//
// $FML: cisco.go,v 1.14 2023/12/30 08:24:17 fukachan Exp $
// $Revision: 1.14 $
//        NAME: cisco.go
// DESCRIPTION: Cisco WiFi device specific module.
//              We suppose only IOS version 12.3.
//              This module configures the wifi device
//		following the given SSID and PASS.

package main

import (
	"fmt"
	"go.bug.st/serial"
	"os"
	"strings"
	"time"
)

type CONFIG [128]string

var count int = 0

// XXX Cisco IOS 12.3 specific configuration
// XXX You must need to prepare each IOS specific initializer.
func ciscoIOS123Init() CONFIG {
	// -------------------- config template begin --------------------
	template := `
!
dot11 ssid {{.SSID}}
   authentication open
   authentication key-management wpa
   guest-mode
   wpa-psk ascii 0 {{.PASSWORD}}
!
interface Dot11Radio0
 no shut
 !
 encryption mode ciphers aes-ccm tkip
 !
 ssid {{.SSID}}
 !
exit
!
interface Dot11Radio1
 shutdown
exit
!
interface FastEthernet0
 no shut
exit
!
no ip http server
no ip http secure-server
!
! exit from "configure" mode
exit
`
	// -------------------- config template end --------------------

	ssid, pass := ConfigGet()
	config := CONFIG{}
	if len(ssid) == 0 || len(pass) == 0 {
		fmt.Fprintf(os.Stderr, "ciscoIOS123Init: empty ssid/pass, ignored\n")
	} else {
		// build the config template with the given SSID and PASS
		config = configBuild(template, ssid, pass)
	}

	// Only if both SSID and PASS is given,
	// return the IOS configuration after SSID and PASS replacement
	return config
}

// build the IOS configuration based on the given template with ssid and pass replacement.
// XXX The replacment is done by "strings" not "text/template" module.
func configBuild(template string, ssid string, pass string) CONFIG {
	// build the config command chain as "config" array
	config := CONFIG{}

	list := strings.Split(template, "\n")
	for i, v := range list {
		w := strings.Replace(v, "{{.SSID}}", ssid, -1)
		w = strings.Replace(w, "{{.PASSWORD}}", pass, -1)
		// be empty for the comment line
		if strings.Contains(w, "!") {
			w = ""
		}

		config[i] = w
		i = i + 1
	}

	return config
}

// We suppose the target device is an initialized Cisco WiFi.
// Hence the prompt must be "ap>" or "ap#" (enabled mode).
// Input data contains config, log and debug messages from Cisco devices,
// so we accept and ignore log and debug messages.
func doConfigure(port serial.Port, data string, isChanged *int) {

	fmt.Fprintf(os.Stderr, "\n(debug) doConfigure> {%s} <IN>\n", data)

	if len(data) == 0 {
		portWrite(port, "\r\n")
		// reset where the initializaion after power off/on must end.
	} else if strings.HasPrefix(data, "Product/Model Number") {
		*isChanged = 1
		ConfigReset()
	} else if strings.HasPrefix(data, "*") {
		fmt.Fprintf(os.Stderr, "\n(debug) doConfigure> {%s} (*)\n", data)
	} else if strings.HasPrefix(data, "%") {
		fmt.Fprintf(os.Stderr, "\n(debug) doConfigure> {%s} (%%)\n", data)
	} else if strings.HasPrefix(data, "^") {
		fmt.Fprintf(os.Stderr, "\n(debug) doConfigure> {%s} (^)\n", data)
	} else if strings.HasPrefix(data, "ap>") {
		portWrite(port, "enable")
	} else if strings.HasPrefix(data, "Password:") {
		portWrite(port, "Cisco")
		portWrite(port, "\r\n")
	} else if strings.HasPrefix(data, "ap#") { // enabled (admin) mode
		portWrite(port, "configure terminal")
	} else if strings.HasPrefix(data, "ap(config)#") { // enabled (admin) configuration mode
		// portWrite(port, "configure terminal")
		// fmt.Fprintf(os.Stderr, "! in configure mode\n")
		configure(port, isChanged)
	}

}

// actually write the IOS command sequence to the target device via serial port
func configure(port serial.Port, isChanged *int) {
	// assert
	if *isChanged == 0 {
		fmt.Fprintf(os.Stderr, "configure> config is not changed, sleep 10\n")
		doWait(10)
		return
	}

	// apply each command in the device specific command chain
	config := ciscoIOS123Init()
	for i, buf := range config {
		if len(buf) > 0 {
			fmt.Printf(" <IN %03d %s\n", i, buf)
		}

		// skip empty lines e.g. next if $buf =~ /^\s*$/;
		strings.Replace(buf, " ", "", -1)
		if len(buf) == 0 {
			continue
		}

		portWrite(port, buf)

		doWait(1)
	}

	*isChanged = 0
}

func doWait(n int) {
	time.Sleep(time.Duration(n) * time.Second)
}
