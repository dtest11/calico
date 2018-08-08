// Copyright (c) 2018 Tigera, Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/colabsaumoh/proto-udsuspver/binder"
	udsver "github.com/colabsaumoh/proto-udsuspver/protos/udsver_v1"
	wlapi "github.com/colabsaumoh/proto-udsuspver/workloadapi"
)

const (
	WorkloadApiUdsHome string = "/tmp/nodeagent"
)

var (
	CfgWldApiUdsHome string

	RootCmd = &cobra.Command{
		Use:   "nodeagent",
		Short: "Node agent with workload api interfaces.",
		Long:  "Node agent with workload api interfaces.",
	}
)

func init() {
	RootCmd.PersistentFlags().StringVarP(&CfgWldApiUdsHome, "wldpath", "w", WorkloadApiUdsHome, "Workload API home path")
}

func Run() {
	// initialize the workload api service
	wl := wlapi.NewWlAPIServer()

	// Create the binder
	b := binder.NewBinder(WorkloadApiUdsHome)

	// Register our service
	udsver.RegisterVerifyServer(b.Server(), wl)

	// Register for system signals
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)

	// Start the binder creating sockets
	bstop := make(chan interface{})
	go b.SearchAndBind(bstop)

	// Wait for term signal.
	<-sigc

	// Shut down the binder.
	bstop <- nil
}

func main() {
	if err := RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}

	// Check if the base directory exisits
	_, e := os.Stat(WorkloadApiUdsHome)
	if e != nil {
		log.Fatalf("WorkloadApi Directory not present (%v)", WorkloadApiUdsHome)
	}

	Run()
}
