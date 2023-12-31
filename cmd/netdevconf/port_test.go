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
// $FML: port_test.go,v 1.3 2023/12/30 10:07:30 fukachan Exp $
// $Revision: 1.3 $
//        NAME: port_test.go
// DESCRIPTION: test routine for USB serial functions.
//

package main

import (
	"fmt"
	"testing"
)

// basic test of serial port probing.
func TestPort_Probe(t *testing.T) {
	port, error := portProbe()

	if error == nil {
		fmt.Printf("port = %s, error = %v\n", port, error)
	} else {
		t.Error("failed to probe USB devices")
	}
}
