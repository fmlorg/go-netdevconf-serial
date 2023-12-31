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
// $FML: db-azure_test.go,v 1.4 2023/12/30 10:07:29 fukachan Exp $
// $Revision: 1.4 $
//        NAME: db-azure_test.go
// DESCRIPTION: check if we can handle Azure Database
//

package main

import (
	"fmt"
	"testing"
)

// import tokens from environmental variables
func TestAzure_readenv(t *testing.T) {
	conn, table := readEnvVar()

	fmt.Printf("endpoint = %s\ntable = %s\n", conn, table)

	if conn != "" && table != "" {
		_ = "ok"
	} else {
		t.Error("failed to read envvar")
	}
}

// check if we can fetch data from Azure
func TestAzure_fetch(t *testing.T) {
	var w WifiConfig

	conn, table := readEnvVar()
	fetch(&w, conn, table)

	fmt.Printf("ssid = %s\npass = %s\n", w.ssid, w.pass)
	if w.ssid != "" && w.pass != "" {
		_ = "ok"
	} else {
		t.Error("failed to fetch from azure-db")
	}
}
