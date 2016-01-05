package config

import (
	"fmt"
	"os"
	//"github.com/davecgh/go-spew/spew"
	"github.com/kelseyhightower/envconfig"
	"github.com/sethdmoore/digo/globals"
	"github.com/sethdmoore/digo/types"
	"strings"
)

func Init() *types.Config {
	var c types.Config
	prefix := strings.ToUpper(globals.APP_NAME)
	err := envconfig.Process(prefix, &c)
	if err != nil {
		fmt.Println("ERROR")
	}

	if c.Trigger == "" {
		c.Trigger = "/bot"
		fmt.Printf("No trigger defined, defaulting to %s\n", c.Trigger)
	}

	if c.ApiInterface == "" {
		c.ApiInterface = "127.0.0.1:8086"
		fmt.Printf("WARN: No API interface defined, defaulting to %s\n", c.ApiInterface)
	}

	if c.Guild == "" {
		fmt.Printf("ERROR: No guild specified to connect to\n")
		fmt.Printf("Digo only connects to a single guild (for now)\n")
		fmt.Printf("Please export %s_SERVER_ID=######\n", prefix)
		fmt.Printf("https://support.discordapp.com/hc/en-us/articles/206346498\n")
		os.Exit(2)
	}

	//spew.Dump(c)
	return &c
}
