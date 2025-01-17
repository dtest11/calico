package main

import (
	"flag"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/dtest11/calico/confd/pkg/buildinfo"
	"github.com/dtest11/calico/confd/pkg/config"
	"github.com/dtest11/calico/confd/pkg/run"
	"github.com/dtest11/calico/libcalico-go/lib/seedrng"
)

var (
	printVersion bool
)

func main() {
	// Make sure the RNG is seeded.
	seedrng.EnsureSeeded()

	flag.BoolVar(&printVersion, "version", false, "print version and exit")
	flag.Parse()
	if printVersion {
		fmt.Printf("confd %s\n", buildinfo.GitVersion)
		os.Exit(0)
	}

	c, err := config.InitConfig(false)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Infof("Config: %#v", c)

	// Run confd.
	run.Run(c)
}
