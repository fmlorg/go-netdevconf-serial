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
// $FML: cisco_test.go,v 1.3 2023/12/30 08:24:17 fukachan Exp $
// $Revision: 1.3 $
//        NAME: cisco_test.go
// DESCRIPTION: test cisto.go
//

package main

import "testing"

// basic test of config loading
func TestCisco_configure(t *testing.T) {
	template := "conf t"

	config := configBuild(template, "test", "test")
	if config[0] == "conf t" && config[1] == "" {
		_ = "ok"
	} else {
		t.Error("failed to configure")
	}
}
