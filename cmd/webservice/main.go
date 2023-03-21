package main

import (
	"flag"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/NoireHub/NATS-streaming-WebService/internal/app/webservice"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath,"","configs/webservice.toml","path to config file")
}

func handleCrash() {
	
}

func main() {
	flag.Parse()

	config:= webservice.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	} 

	if err := webservice.Start(config); err != nil {
		log.Fatal(err)
	}
}