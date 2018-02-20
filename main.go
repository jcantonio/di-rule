package main

import (
	"fmt"

	"github.com/jcantonio/di-rule/api"
	"github.com/jcantonio/di-rule/command"
	"github.com/jinzhu/configor"
)

var Config = struct {
	DB struct {
		Name     string `default:"di-rule"`
		User     string `default:""`
		Password string `default:""`
		Address  string `default:"http://localhost"`
		Port     uint   `default:"5984"`
	}
	Server struct {
		Port uint `default:"8000"`
	}
}{}

func main() {
	err := configor.Load(&Config, "config.yml")
	if err != nil {
		fmt.Printf("Init DB %s", err)
	}
	dbURL := fmt.Sprintf("%s:%d", Config.DB.Address, Config.DB.Port)
	fmt.Printf("Init DB %s \n", dbURL)
	command.InitDatabase(dbURL, Config.DB.Name)
	command.LoadRulesInMem()
	fmt.Printf("Start server on port %d\n", Config.Server.Port)
	api.Init(Config.Server.Port)
}
