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
// $FML: config.go,v 1.3 2023/12/28 14:05:03 fukachan Exp $
// $Revision: 1.3 $
//        NAME: config.go
// DESCRIPTION: hold and manipulate Wifi's most important variables e.g. SSID and PASS.
//

package main

type WifiConfig struct {
	ssid string
	pass string
}

var config WifiConfig
var configold WifiConfig

func ConfigGet() (string, string) {
	return config.ssid, config.pass
}

func ConfigSet(ssid string, pass string) {

	// backup the old config
	configold.ssid = config.ssid
	configold.pass = config.pass

	// update the current config
	config.ssid = ssid
	config.pass = pass
}

func ConfigIsSame(ssid string, pass string) bool {

	if config.ssid == ssid && config.pass == pass {
		return true
	} else {
		return false
	}
}

func ConfigReset() {
	ConfigSet("", "")
}
