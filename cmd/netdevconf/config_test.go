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
// $FML: config_test.go,v 1.3 2023/12/28 14:05:04 fukachan Exp $
// $Revision: 1.3 $
//        NAME: config_test.go
// DESCRIPTION: config.go test rules
//

package main

import "testing"

func TestConfig1(t *testing.T) {
	var a WifiConfig
	var b WifiConfig
	a.ssid = "testssid"
	a.pass = "testpass"

	ConfigSet("testssid", "testpass")
	b.ssid, b.pass = ConfigGet()

	if a.ssid == b.ssid && a.pass == b.pass {
		_ = "ok" // noop
	} else {
		t.Error("not updated 1")
	}
}

func TestConfigIsSame_same(t *testing.T) {

	ConfigSet("testssid2", "testpass2")
	if ConfigIsSame("testssid2", "testpass2") {
		_ = "ok" // noop
	} else {
		t.Error("must be same")
	}
}

func TestConfigIsSame_differ(t *testing.T) {

	ConfigSet("testssid1", "testpass1")
	if ConfigIsSame("testssid2", "testpass2") {
		t.Error("must be different")
	} else {
		_ = "ok" // noop
	}
}

func TestConfig_reset(t *testing.T) {

	ConfigSet("testssid1", "testpass1")

	ConfigReset()

	ssid, pass := ConfigGet()
	if ssid == "" && pass == "" {
		_ = "ok" // noop
	} else {
		t.Error("must be empty")
	}
}
