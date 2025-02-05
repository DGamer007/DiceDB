// Copyright (c) 2022-present, DiceDB contributors
// All rights reserved. Licensed under the BSD 3-Clause License. See LICENSE file in the project root for full license information.

package ironhawk

import (
	"os"
	"sync"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	var wg sync.WaitGroup
	// Run the test server
	// This is a synchronous method, because internally it
	// checks for available port and then forks a goroutine
	// to start the server
	opts := TestServerOptions{
		Port: 9739,
	}
	RunTestServer(&wg, opts)

	// Wait for the server to start
	time.Sleep(2 * time.Second)

	// Run the test suite
	exitCode := m.Run()

	client := getLocalConnection()
	if client == nil {
		panic("Failed to connect to the test server")
	}
	defer client.Close()
	result := client.FireString("ABORT")
	if result.GetVStr() != "OK" {
		panic("Failed to abort the server")
	}

	wg.Wait()
	os.Exit(exitCode)
}
