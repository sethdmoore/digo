package config

import (
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Email    string `envconfig:"user" json:"email"`
	Password string `envconfig:"pass" json:"password"`
	Trigger  string `envconfig:"trigger" json:"trigger"`
	//Token    string `json:"token"`
}

func Init() *Config {
	var c Config
	err := envconfig.Process("discordio", &c)
	if err != nil {
		fmt.Println("ERROR")
	}
	if c.Trigger == "" {
		fmt.Println("No trigger defined, defaulting to /bot")
		c.Trigger = "/bot"
	}

	// pad the trigger with a space
	c.Trigger += " "
	//spew.Dump(c)
	return &c
}
