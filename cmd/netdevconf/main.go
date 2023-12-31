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
// $FML: main.go,v 1.5 2023/12/30 08:31:58 fukachan Exp $
// $Revision: 1.5 $
//        NAME: main.go
// DESCRIPTION: the main function of a tool to configure network devices e.g. Cisco WiFi.
//              We suppose Cisco WiFi-AP running IOS 12.3, connected to the PC via
//              PL2303 based USB serial cable.
//

package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	timeOut := 3   // 3 seconds
	isChanged := 1 // do it in the boot time, check if the cloud data is changed

	//create go channels
	dataStream := make(chan string, 1)
	azureStream := make(chan int, 1)
	readerr := make(chan error)

	// open USB serial port
	port := portOpen()

	// IO loop (preparation)
	// - run a go channel to read the serial port
	// - run another go channel to connect the Azure Database
	go portRead(port, dataStream, readerr)
	go azureFetchConfig(azureStream)

	// IO loop (main; start the actual IO loop)
	portWrite(port, "")
loop:
	for {
		select {
		case data := <-dataStream:
			doConfigure(port, data, &isChanged)
		case <-azureStream:
			isChanged = 1
		case <-time.After(time.Duration(timeOut) * time.Second): // 3sec.
			fmt.Fprintln(os.Stderr, "timeout")
			portWrite(port, "")
		case <-time.After(30 * time.Duration(timeOut) * time.Second): // 30x3 = 90sec.
			fmt.Fprintln(os.Stderr, "enforce re-configure after 30*timeout")
			// enforce reset each 90sec. where 90 is a magic number under no deep meaning.
			isChanged = 1
			ConfigReset()
		case err := <-readerr:
			if err != nil {
				log.Fatal(err)
			}
			break loop
		}
	}
}
