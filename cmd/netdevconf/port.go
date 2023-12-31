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
// $FML: port.go,v 1.7 2023/12/30 10:07:30 fukachan Exp $
// $Revision: 1.7 $
//        NAME: port.go
// DESCRIPTION: handle USB serial port which chip is supposed to be PL2303.
//

package main

import (
	"bufio"
	"errors"
	"fmt"
	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
	"log"
	"os"
)

// probe USB serial port to find PL2303.
// This prober makes the code OS-independent.
// Currently we support only PL2303 chip.
func portProbe() (string, error) {
	is_found := 0
	port_found := ""

	ports, error := enumerator.GetDetailedPortsList()
	if error != nil {
		return "", error
	}
	for _, port := range ports {
		if port.IsUSB {
			fmt.Printf("(debug) FOUND DEVICE %s %s %s\n", port.Name, port.VID, port.PID)
		}

		// PL2303 Vender=067b Product=2303 (067b (UNIX), 067B (Win11))
		if port.IsUSB && (port.VID == "067b" || port.VID == "067B") && port.PID == "2303" {
			fmt.Printf("(debug) FOUND PL2303 %s %s %s\n", port.Name, port.VID, port.PID)
			is_found = is_found + 1
			port_found = port.Name
		}
	}

	if is_found > 0 {
		return port_found, nil
	} else {
		return "", errors.New("no serial")
	}
}

// open serial port in the mode "9600 8N1 9600 8bit no-parity 1-stopbit"
// which setting is typical for traditional network devices and PCs (D-Sub 9pins, RS-232C).
func portOpen() serial.Port {
	var bits serial.ModemOutputBits
	bits.RTS = false
	bits.DTR = false

	mode := &serial.Mode{
		BaudRate:          9600,
		DataBits:          8,
		Parity:            serial.NoParity,
		StopBits:          serial.OneStopBit,
		InitialStatusBits: &bits,
	}

	portName, err := portProbe()
	if err != nil {
		log.Fatal(err)
	}

	port, err := serial.Open(portName, mode)
	if err != nil {
		log.Fatal(err)
	}

	return port
}

func portWrite(port serial.Port, data string) {
	fmt.Fprintf(os.Stderr, "write <%s>\n", data)

	n, err := port.Write([]byte(data + "\r\n"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(os.Stderr, "Sent %v bytes\n", n)
}

// send the read data to the go channel which main() waits for.
// It is a trampoline: portRead() -> main() -> doConfigure() -> portWrite()
func portRead(port serial.Port, data chan string, dataerr chan error) {

	scanner := bufio.NewScanner(port)

	for scanner.Scan() {
		data <- scanner.Text()
	}
	close(data) // close causes the range on the channel to break out of the loop
	dataerr <- scanner.Err()
}
