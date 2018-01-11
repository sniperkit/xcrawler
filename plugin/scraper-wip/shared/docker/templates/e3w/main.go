package main

import (
	"flag"
	"fmt"
	"os"

	// "github.com/k0kubun/pp"
	"github.com/roscopecoltran/configor"

	"github.com/coreos/etcd/version"
	"github.com/gin-gonic/gin"

	"github.com/roscopecoltran/e3w/conf"
	"github.com/roscopecoltran/e3w/e3ch"
	"github.com/roscopecoltran/e3w/routers"
)

const (
	PROGRAM_NAME    = "e3w"
	PROGRAM_VERSION = "0.0.2"
)

var (
	configFilepath  string
	frontPrefixPath string
	useConfigor     bool
)

func init() {
	flag.StringVar(&configFilepath, "conf", "/data/conf.d/e3w/config.ini", "config file path")
	flag.StringVar(&frontPrefixPath, "front-dir", "/data/static/dist", "frontend static dir path")
	flag.BoolVar(&useConfigor, "use-configor", false, "use configor to load app config (supported formats: *.json, *.yaml or *.toml)")

	rev := flag.Bool("rev", false, "print rev")
	flag.Parse()

	if *rev {
		fmt.Printf("[%s v%s]\n[etcd %s]\n",
			PROGRAM_NAME, PROGRAM_VERSION,
			version.Version,
		)
		os.Exit(0)
	}
}

func main() {
	config, err := conf.Init(configFilepath)
	if err != nil {
		panic(err)
	}

	// pp.Print(config)
	if frontPrefixPath != "" {
		config.Front.Dist.Dir = frontPrefixPath
	}

	if err := configor.Dump(config, "e3w", "yaml,json,toml", "./shared/conf.d/e3w/dumps"); err != nil {
		panic(err)
	}

	client, err := e3ch.NewE3chClient(config)
	if err != nil {
		panic(err)
	}

	router := gin.Default()
	routers.InitRouters(router, config, client)
	router.Run(":" + config.Port)

}
