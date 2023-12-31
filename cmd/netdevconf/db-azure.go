// -*- utf-8 -*-
//
// Copyright (C) 2023 Shunsuke Toyosaki
// Copyright (C) 2023 Ken'ichi Fukamachi
//   All rights reserved. This program is free software; you can
//   redistribute it and/or modify it under 2-Clause BSD License.
//   https://opensource.org/licenses/BSD-2-Clause
//
// mailto: fukachan@fml.org
//    web: https://www.fml.org/
// github: https://github.com/fmlorg
//
// $FML: db-azure.go,v 1.10 2023/12/30 10:07:29 fukachan Exp $
// $Revision: 1.10 $
//        NAME: db-azure.go
// DESCRIPTION: fetch data from Azure Database.
//              There are top level loop and utitity functions.
//

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
)

// run this function infinitely under the go channel.
// - check if the configuration has changed or not
// - send a signal (we need to reconfigure the device) via go channel if changed.
//   XXX the configuration is saved in config.go via ConfigSet() (see "config.go").
func azureFetchConfig(c chan int) {
	var w WifiConfig

	for {
		time.Sleep(30 * time.Second)

		conn, table := readEnvVar()
		fetch(&w, conn, table)

		if ConfigIsSame(w.ssid, w.pass) {
			fmt.Fprintf(os.Stderr, "azureFetchConfig: ssid,pass same\n")
		} else {
			fmt.Fprintf(os.Stderr, "azureFetchConfig: ssid,pass updated\n")
			ConfigSet(w.ssid, w.pass)
			c <- 1
		}
	}

}

func readEnvVar() (string, string) {
	connStr := os.Getenv("AZURE_ENDPOINT")
	if connStr == "" {
		panic("specify the mandatory environemntal variable AZURE_ENDPOINT")
	}
	tableName := os.Getenv("AZURE_TABLE")
	if connStr == "" {
		panic("specify the mandatory environemntal variable AZURE_TABLE")
	}

	return connStr, tableName
}

// fetch the data from Azure Database and return the latest one
// where we assume the azure data
//
func fetch(w *WifiConfig, connStr string, tableName string) {
	todayDate := time.Now().UTC()
	type WifiEntity struct {
		aztables.Entity
		SSID     string
		Password string
	}

	// Initialize a slice to hold matched wifiEntity entries.
	// This slice is used to sort to find the latest entry.
	x := []WifiEntity{}

	// try to connect to Azure
	serviceClient, err := aztables.NewServiceClientFromConnectionString(connStr, nil)
	if err != nil {
		panic(err)
	}

	// Next, after the connection to azure has been established, init the list pager for the table.
	// NewListEntitiesPager(nil) means "use the default options".
	listPager := serviceClient.NewClient(tableName).NewListEntitiesPager(nil)

	var pageCount int = 0
	for listPager.More() {
		response, err := listPager.NextPage(context.TODO())
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(os.Stderr, "There are %d entities in page #%d\n", len(response.Entities), pageCount)

		pageCount += 1
		for _, entity := range response.Entities {
			var wifiEntity WifiEntity
			err = json.Unmarshal(entity, &wifiEntity)
			if err != nil {
				panic(err)
			}

			// append the entry if its YYYYMM is same as today's YYYYMM.
			if (todayDate.Year() == time.Time(wifiEntity.Timestamp).Year()) && (todayDate.Month() == time.Time(wifiEntity.Timestamp).Month()) {
				x = append(x, wifiEntity)
			}

			// debug, always ok
			x = append(x, wifiEntity)
		}
	}

	fmt.Fprintf(os.Stderr, "%v\n", x)

	// sort the slice in timeorder
	sort.Slice(x, func(i, j int) bool { return time.Time(x[i].Timestamp).Before(time.Time(x[j].Timestamp)) })

	// fmt.Fprintf(os.Stderr, "%x\n", x)

	// get the latest entry
	ssid := ""
	pass := ""
	if len(x) > 0 {
		ssid = x[len(x)-1].SSID
		pass = x[len(x)-1].Password
	}

	// debug ?
	fmt.Fprintf(os.Stdout, "%s, %s\n", ssid, pass)

	// save the fetched data in &w to return it to azureFetchConfig()
	w.ssid = ssid
	w.pass = pass
}
