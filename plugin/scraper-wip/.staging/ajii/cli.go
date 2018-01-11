package ajii

import (
	"fmt"
	"gopkg.in/urfave/cli.v1"
	"log"
)

func startService(config *EtcdConfig, service Service) interface{} {
	return func(c *cli.Context) error {
		// servicename := c.String("name")
		// if servicename == "" {
		// 	log.Fatal("You must to set a servicename")
		// }
		// service := NewService(config, servicename)
		service.Init(c)
		return nil
	}
}

func dumpConfig(config *EtcdConfig) interface{} {
	return func(c *cli.Context) error {
		nodes, _ := config.Dump()
		idx := len(nodes)
		if idx > 0 {
			for _, node := range nodes {
				fmt.Printf("%s    : %s\n", node.Key, node.Value)
			}
		} else {
			fmt.Println("No keys in store")
		}
		return nil
	}
}

func getConfig(config *EtcdConfig) interface{} {
	return func(c *cli.Context) error {
		if c.NArg() > 0 {
			for _, key := range c.Args() {
				node, _ := config.Get(key)
				if node.Value != "" {
					fmt.Printf("%s    : %s\n", node.Key, node.Value)
				} else {
					fmt.Printf("Key `%s` not found\n", key)
				}
			}
		} else {
			log.Fatal("You must specify a key")
		}
		return nil
	}
}

func setConfig(config *EtcdConfig) interface{} {
	return func(c *cli.Context) error {
		if c.NArg() == 2 {
			key := c.Args().Get(0)
			value := c.Args().Get(1)
			force := c.Bool("force")
			node, _ := config.Get(key)
			if node.Value == "" || force {
				config.Set(key, value)
				fmt.Printf("Key `/%s` set to `%s`\n", key, value)
			} else {
				log.Fatalf("Can not overwrite existing key: `%s` is set to `%s`", node.Key, node.Value)
			}
		} else {
			log.Fatal("You must specify a key and a value")
		}
		return nil
	}
}

func NewCli(config *EtcdConfig, service Service) *cli.App {
	app := cli.NewApp()
	app.Name = "AJII"
	app.Version = "0.1"
	app.Usage = "A Journey Into the Impossible"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "Connection string to kv-store",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:   "run",
			Usage:  "Run a thing",
			Action: startService(config, service),
			Flags: []cli.Flag{
				// cli.StringFlag{
				// 	Name:  "name, n",
				// 	Usage: "name of thing",
				// },
				cli.StringFlag{
					Name:  "host, ho",
					Usage: "hostname to announce",
					Value: "localhost",
				},
				cli.IntFlag{
					Name:  "port, p",
					Usage: "port to bind frontend",
					Value: 8080,
				},
			},
		},
		{
			Name:  "config",
			Usage: "Manage configuration",
			Subcommands: []cli.Command{
				{
					Name:   "dump",
					Usage:  "Dump current configuration",
					Action: dumpConfig(config),
				},
				{
					Name:   "get",
					Usage:  "Get value of specified key(s))",
					Action: getConfig(config),
				},
				{
					Name:  "set",
					Usage: "Set a key with a specified value",
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name:  "force, f",
							Usage: "Force setting the key if it exists",
						},
					},
					Action: setConfig(config),
				},
			},
		},
	}
	return app
}
